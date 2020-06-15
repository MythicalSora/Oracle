package commands

import (
	"errors"

	"github.com/mythicalsora/Oracle/models"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"github.com/mythicalsora/Oracle/libs"
)

var guild *discordgo.Guild

//Housing ...
func Housing(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB, c *discordgo.Channel, args []string) {
	session = s
	message = m
	channel = c
	dbp = db
	guild, _ = session.Guild(message.GuildID)

	if len(args) >= 1 {
		switch args[0] {
		case "list":

		case "sell":

		case "buy":
			var user models.User
			dbp.First(&user, "discord_id = ?", message.Author.ID)

			if user.Balance >= 50000.00 {
				if len(args) > 1 {
					name := args[1] + "'s-house"
					cat := discordgo.GuildChannelCreateData{
						Name:     name,
						Type:     discordgo.ChannelTypeGuildText,
						ParentID: "716938410559799297",
					}

					st, err := session.GuildChannelCreateComplex(guild.ID, cat)

					if err != nil {
						libs.Error(err, session, channel)
					}

					house := models.House{
						Channel: st.ID,
						Price:   50000.00,
						Owner:   user.ID,
					}

					if err := dbp.Create(&house).Error; err != nil {
						libs.Error(err, session, channel)
					} else {
						user.Balance = user.Balance - 50000.00
						dbp.Save(&user)
						session.ChannelMessageSendEmbed(channel.ID, libs.HousingBuyEmbed(message.Author.Username, st))
					}
				}
			}
		}
	} else {
		e := errors.New("please pick one of the following options: ``list`` | ``sell`` | ``buy``")
		libs.Error(e, session, channel)
	}
}
