package vfm

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/miso"
	vault "github.com/curtisnewbie/user-vault/api"
)

const (
	// unpackNotiTitleMax and unpackNotiMessageMax are the limits user-vault
	// validates the notification against.
	unpackNotiTitleMax   = 255
	unpackNotiMessageMax = 1000

	// unpackNotiMoreReserve is room kept for a shortened path list's ", and N
	// more" tail.
	unpackNotiMoreReserve = 16
)

// unpackReport is how one unpack run ended. At most one notification is built
// from it, so a run that both dropped entries and stopped early says both
// things in one message.
type unpackReport struct {
	ArchiveName string // name of the zip file, empty when unknown
	TargetName  string // name of the unpack target directory, empty when unknown
	TotalFiles  int    // number of file entries the archive offered
	Created     int    // number of file entries written
	Drops       []unpackDrop
	Stopped     bool // a row write failed, the remaining entries were not attempted
}

// notify sends the report through user-vault, unless there is nothing to say.
// A notification that cannot be sent is logged and never changes the outcome.
func (r unpackReport) notify(rail miso.Rail, userNo string) {
	if userNo == "" {
		return
	}
	title, message := r.render()
	if title == "" {
		return
	}
	sendUnpackNotification(rail, userNo, title, message)
}

// unpackRefusal returns the error mini-fstore refused the unpack with, or nil
// when the call failed for another reason such as a transport error.
//
// The refusal code sits on the innermost MisoErr, so the chain is walked.
func unpackRefusal(err error) *errs.MisoErr {
	for err != nil {
		if me, ok := err.(*errs.MisoErr); ok && me.HasCode() {
			return me
		}
		err = errors.Unwrap(err)
	}
	return nil
}

// unpackRefusalReason returns what to tell the user about a refused archive, or
// "" when the error is not a refusal.
//
// The reason is mini-fstore's own message, so the limit and the offending entry
// come from the one place that read the archive.
func unpackRefusalReason(err error) string {
	me := unpackRefusal(err)
	if me == nil {
		return ""
	}
	switch me.Code() {
	case fstore.ZipTooDeep, fstore.IllegalFormat:
	default:
		return ""
	}

	reason := strings.TrimRight(strings.TrimSpace(me.Msg()), ". ")
	if reason == "" {
		reason = "the archive was refused"
	}
	return reason
}

// render builds the one message this run earns. Both a stop and dropped entries
// share it, so the user never gets two notifications for one unpack.
func (r unpackReport) render() (title, message string) {
	if !r.Stopped && len(r.Drops) == 0 {
		return "", ""
	}

	archive := unpackQuoted(r.ArchiveName, "the archive")
	target := unpackQuoted(r.TargetName, "the target folder")
	var head, tail strings.Builder

	if r.Stopped {
		title = "Archive unpack stopped early"
		notCreated := r.TotalFiles - r.Created - len(r.Drops)
		if notCreated < 0 {
			notCreated = 0
		}
		fmt.Fprintf(&head, "The archive %s stopped unpacking into %s after an internal error. %d of %d entries were not created",
			archive, target, notCreated, r.TotalFiles)
		tail.WriteString(" The files already written were kept.")
	} else {
		title = "Archive unpacked with skipped entries"
		fmt.Fprintf(&head, "The archive %s was unpacked into %s", archive, target)
	}

	ended := false
	if len(r.Drops) > 0 {
		if r.Stopped {
			head.WriteString(", and")
		} else {
			head.WriteString(", but")
		}
		fmt.Fprintf(&head, " skipped %s: %s", unpackEntryCount(len(r.Drops)), unpackDropGroups(r.Drops))

		// leave room for the paths and the closing sentence
		budget := unpackNotiMessageMax - unpackRuneLen(head.String()) - unpackRuneLen(tail.String()) - unpackRuneLen(". Skipped paths: .") - unpackNotiMoreReserve
		if budget < 80 {
			budget = 80
		}
		fmt.Fprintf(&head, ". Skipped paths: %s.", unpackDropPaths(r.Drops, budget))
		ended = true
	}
	if !ended {
		head.WriteString(".")
	}
	return title, unpackTruncate(head.String()+tail.String(), unpackNotiMessageMax)
}

// unpackDropGroups names each drop reason with how many entries it covers, in
// the order the reasons first appeared. The counts add up to the total stated
// before the list.
func unpackDropGroups(drops []unpackDrop) string {
	type group struct {
		reason string
		count  int
	}
	groups := make([]group, 0, len(drops))
	index := make(map[string]int, len(drops))
	for _, d := range drops {
		i, ok := index[d.Reason]
		if !ok {
			i = len(groups)
			index[d.Reason] = i
			groups = append(groups, group{reason: d.Reason})
		}
		groups[i].count++
	}

	var b strings.Builder
	for i, g := range groups {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%s (%d)", g.reason, g.count)
	}
	return b.String()
}

// unpackEntryCount names a count of entries, so a message never says "1 entries".
func unpackEntryCount(n int) string {
	if n == 1 {
		return "1 entry"
	}
	return fmt.Sprintf("%d entries", n)
}

// unpackDropPaths joins the dropped entry paths for a message, shortening the
// list rather than losing the counts when there are too many of them.
func unpackDropPaths(drops []unpackDrop, budget int) string {
	var b strings.Builder
	used := 0
	for i, d := range drops {
		p := d.Path
		if p == "" {
			p = "(unnamed entry)"
		}
		w := unpackRuneLen(p)
		if used > 0 && used+w+2 > budget {
			fmt.Fprintf(&b, ", and %d more", len(drops)-i)
			break
		}
		if used > 0 {
			b.WriteString(", ")
			used += 2
		}
		b.WriteString(p)
		used += w
	}
	return b.String()
}

// sendUnpackNotification delivers one notification. A lost notification is only
// logged, never a reason to fail the unpack.
func sendUnpackNotification(rail miso.Rail, userNo, title, message string) {
	title = unpackTruncate(title, unpackNotiTitleMax)
	message = unpackTruncate(message, unpackNotiMessageMax)
	if err := vault.CreateNotification(rail, vault.CreateNotificationReq{
		Title:           title,
		Message:         message,
		ReceiverUserNos: []string{userNo},
	}); err != nil {
		rail.Errorf("Failed to create unpack notification, userNo: %v, title: %v, %v", userNo, title, err)
	}
}

// unpackQuoted renders a name for the message, falling back to a description
// when the name is unknown.
func unpackQuoted(name, fallback string) string {
	if name == "" {
		return fallback
	}
	return fmt.Sprintf("%q", name)
}

func unpackRuneLen(s string) int {
	return utf8.RuneCountInString(s)
}

// unpackTruncate shortens a string to max runes, keeping it valid UTF-8.
func unpackTruncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-3]) + "..."
}
