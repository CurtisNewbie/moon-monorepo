package vfm

import (
	"errors"
	"strings"
	"testing"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
)

func TestUnpackRefusal(t *testing.T) {
	deep := errs.NewErrfCode(fstore.ZipTooDeep, "Zip file is nested too deeply, max depth: 32")
	if got := unpackRefusal(errs.Wrapf(deep, "failed to TriggerFileUnZip")); got == nil || got.Code() != fstore.ZipTooDeep {
		t.Fatalf("refusal: got %+v, want code %q", got, fstore.ZipTooDeep)
	}
	if got := unpackRefusal(errs.NewErrf("no code here")); got != nil {
		t.Fatalf("refusal: got %+v, want none", got)
	}
}

func TestUnpackRefusalReason(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "depth refusal carries mini-fstore's own message",
			err:  errs.NewErrfCode(fstore.ZipTooDeep, "Zip file is nested too deeply, max depth: 32, deepest entry: a/b/c"),
			want: "Zip file is nested too deeply, max depth: 32, deepest entry: a/b/c",
		},
		{
			name: "unreadable archive is a refusal too",
			err:  errs.NewErrfCode(fstore.IllegalFormat, "Illegal zip file format"),
			want: "Illegal zip file format",
		},
		{
			name: "a transport failure is not a refusal",
			err:  errors.New("dial tcp 127.0.0.1:8084: connect: connection refused"),
			want: "",
		},
		{
			name: "a coded error outside the refusal codes is not a refusal",
			err:  errs.NewErrfCode(fstore.FileNotFound, "file not found"),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unpackRefusalReason(tt.err); got != tt.want {
				t.Fatalf("reason: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnpackTruncate(t *testing.T) {
	long := strings.Repeat("é", unpackNotiMessageMax+10)
	got := unpackTruncate(long, unpackNotiMessageMax)
	if len([]rune(got)) != unpackNotiMessageMax {
		t.Fatalf("truncated runes: got %d, want %d", len([]rune(got)), unpackNotiMessageMax)
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("truncated message does not end with a marker: %q", got[len(got)-10:])
	}
	if strings.Count(got, "é") < unpackNotiMessageMax-5 {
		t.Fatalf("truncation cut more than it had to")
	}
}

func TestUnpackReportRender(t *testing.T) {
	tests := []struct {
		name         string
		report       unpackReport
		wantTitle    string
		wantInMsg    []string
		wantNotInMsg []string
	}{
		{
			name: "drops only",
			report: unpackReport{
				ArchiveName: "photos.zip",
				TargetName:  "photos.zip unpacked 20260101_120000",
				TotalFiles:  10,
				Created:     8,
				Drops: []unpackDrop{
					{Path: "../evil.txt", Reason: "the path contains a '..' segment"},
					{Path: "a/b.", Reason: "a name ends with a dot or a space"},
				},
			},
			wantTitle: "Archive unpacked with skipped entries",
			wantInMsg: []string{
				`"photos.zip"`, `"photos.zip unpacked 20260101_120000"`,
				"skipped 2 entries",
				"the path contains a '..' segment (1)",
				"a name ends with a dot or a space (1)",
				"Skipped paths: ../evil.txt, a/b.",
			},
			wantNotInMsg: []string{"partial", "not created"},
		},
		{
			name: "reasons are grouped in first seen order",
			report: unpackReport{
				ArchiveName: "grouped.zip",
				TotalFiles:  10,
				Created:     6,
				Drops: []unpackDrop{
					{Path: "a.", Reason: "a name ends with a dot or a space"},
					{Path: "../b", Reason: "the path contains a '..' segment"},
					{Path: "c ", Reason: "a name ends with a dot or a space"},
					{Path: "d.", Reason: "a name ends with a dot or a space"},
				},
			},
			wantTitle: "Archive unpacked with skipped entries",
			wantInMsg: []string{
				"skipped 4 entries",
				"a name ends with a dot or a space (3)",
				"the path contains a '..' segment (1)",
			},
		},
		{
			name: "stopped early",
			report: unpackReport{
				ArchiveName: "big.zip",
				TotalFiles:  40,
				Created:     12,
				Stopped:     true,
			},
			wantTitle:    "Archive unpack stopped early",
			wantInMsg:    []string{`"big.zip"`, "after an internal error", "28 of 40 entries were not created", "already written were kept"},
			wantNotInMsg: []string{"storage", "Skipped paths"},
		},
		{
			name: "one message covers both",
			report: unpackReport{
				ArchiveName: "mixed.zip",
				TotalFiles:  10,
				Created:     4,
				Stopped:     true,
				Drops:       []unpackDrop{{Path: "../evil.txt", Reason: "the path contains a '..' segment"}},
			},
			wantTitle: "Archive unpack stopped early",
			wantInMsg: []string{
				"5 of 10 entries were not created",
				"and skipped 1 entry: the path contains a '..' segment (1)",
				"Skipped paths: ../evil.txt",
				"already written were kept",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, msg := tt.report.render()
			if title != tt.wantTitle {
				t.Fatalf("title: got %q, want %q", title, tt.wantTitle)
			}
			if len([]rune(title)) > unpackNotiTitleMax {
				t.Fatalf("title longer than %d: %d", unpackNotiTitleMax, len([]rune(title)))
			}
			if len([]rune(msg)) > unpackNotiMessageMax {
				t.Fatalf("message longer than %d: %d", unpackNotiMessageMax, len([]rune(msg)))
			}
			for _, want := range tt.wantInMsg {
				if !strings.Contains(msg, want) {
					t.Fatalf("message %q does not contain %q", msg, want)
				}
			}
			for _, not := range tt.wantNotInMsg {
				if strings.Contains(msg, not) {
					t.Fatalf("message %q contains %q", msg, not)
				}
			}
		})
	}
}

func TestUnpackReportRenderShortensALongDropList(t *testing.T) {
	drops := make([]unpackDrop, 0, 200)
	for i := 0; i < 200; i++ {
		drops = append(drops, unpackDrop{Path: strings.Repeat("p", 40), Reason: "the path contains a '..' segment"})
	}
	report := unpackReport{ArchiveName: "big.zip", TotalFiles: 300, Created: 100, Drops: drops}
	_, msg := report.render()

	if len([]rune(msg)) > unpackNotiMessageMax {
		t.Fatalf("message not shortened, runes: %d", len([]rune(msg)))
	}
	if !strings.Contains(msg, "skipped 200 entries") {
		t.Fatalf("message lost the drop count: %q", msg)
	}
	if !strings.Contains(msg, "the path contains a '..' segment (200)") {
		t.Fatalf("message lost the grouped reason: %q", msg)
	}
	if !strings.Contains(msg, "more") {
		t.Fatalf("message lost the shortened list marker: %q", msg)
	}
}

func TestUnpackReportRenderStaysWithinTheLimitWithManyReasons(t *testing.T) {
	drops := make([]unpackDrop, 0, 300)
	for i := 0; i < 300; i++ {
		drops = append(drops, unpackDrop{
			Path:   strings.Repeat("q", 30),
			Reason: "reason " + strings.Repeat("r", 20) + string(rune('a'+i%20)),
		})
	}
	report := unpackReport{ArchiveName: "messy.zip", TotalFiles: 300, Created: 0, Stopped: true, Drops: drops}
	_, msg := report.render()

	if len([]rune(msg)) > unpackNotiMessageMax {
		t.Fatalf("message not shortened, runes: %d", len([]rune(msg)))
	}
	if !strings.Contains(msg, "skipped 300 entries") {
		t.Fatalf("message lost the drop count: %q", msg)
	}
}

func TestUnpackReportRenderSaysNothingWhenNothingHappened(t *testing.T) {
	report := unpackReport{ArchiveName: "fine.zip", TotalFiles: 3, Created: 3}
	if title, msg := report.render(); title != "" || msg != "" {
		t.Fatalf("a clean run still built a message: %q %q", title, msg)
	}
}
