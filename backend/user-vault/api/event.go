package api

import "github.com/curtisnewbie/miso/middleware/rabbit"

var (
	CreateNotifiPipeline = rabbit.NewEventPipeline[CreateNotifiEvent]("pieline.user-vault.create-notifi").
				LogPayload().
				MaxRetry(3).
				Document("CreateNotifiPipeline", "Pipeline that creates notifications to the specified list of users", "user-vault")

	CreateNotifiByAccessPipeline = rabbit.NewEventPipeline[CreateNotifiByAccessEvent]("pieline.user-vault.create-notifi.by-access").
					LogPayload().
					MaxRetry(3).
					Document("CreateNotifiByAccessPipeline", "Pipeline that creates notifications to users who have access to the specified resource", "user-vault")
)

type CreateNotifiEvent struct {
	Title           string   `valid:"maxLen:255" desc:"notification title"`
	Message         string   `valid:"maxLen:1000" desc:"notification content"`
	ReceiverUserNos []string `desc:"user_no of receivers"`
}

type CreateNotifiByAccessEvent struct {
	Title   string `valid:"maxLen:255" desc:"notification title"`
	Message string `valid:"maxLen:1000" desc:"notification content"`
	ResCode string `valid:"notEmpty" desc:"resource code"`
}
