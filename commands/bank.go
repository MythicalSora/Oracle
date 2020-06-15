package commands

import (
	"strconv"
	"time"

	"github.com/mythicalsora/Oracle/libs"
	"github.com/mythicalsora/Oracle/models"

	"github.com/bwmarrin/discordgo"
	// Justification ...
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var dbp *gorm.DB
var session *discordgo.Session
var channel *discordgo.Channel
var message *discordgo.MessageCreate

// Bank ...
func Bank(s *discordgo.Session, m *discordgo.MessageCreate, db *gorm.DB, c *discordgo.Channel, args []string) {
	dbp = db
	session = s
	channel = c
	message = m

	if len(args) >= 1 {
		switch args[0] {
		case "add":
			guild, _ := session.Guild(message.GuildID)
			if libs.CheckUser(session, guild, message) == true {
				amount, err := strconv.ParseFloat(args[2], 32)
				if err == nil == true {
					var target models.User
					dbp.First(&target, "discord_id = ?", message.Mentions[0].ID)

					target.Balance = target.Balance + amount
					dbp.Save(&target)

					session.ChannelMessageSendEmbed(channel.ID, libs.AddEmbed(message.Mentions[0].Username))
				}
			}

		case "sub":
			guild, _ := session.Guild(message.GuildID)
			if libs.CheckUser(session, guild, message) == true {
				amount, err := strconv.ParseFloat(args[2], 32)
				if err == nil == true {
					var target models.User
					dbp.First(&target, "discord_id = ?", message.Mentions[0].ID)

					target.Balance = target.Balance - amount
					dbp.Save(&target)

					session.ChannelMessageSendEmbed(channel.ID, libs.SubEmbed(message.Mentions[0].Username))
				}
			}
		case "pay":
			amount, err := strconv.ParseFloat(args[2], 32)

			if err == nil {
				if checkInput(amount) == true {
					var user, target models.User
					dbp.First(&user, "discord_id = ?", message.Author.ID)
					dbp.First(&target, "discord_id = ?", message.Mentions[0].ID)

					user.Balance = user.Balance - amount
					target.Balance = target.Balance + amount

					dbp.Save(&user)
					dbp.Save(&target)

					session.ChannelMessageSendEmbed(channel.ID, libs.BankPayEmbed(message.Mentions[0].Username, amount))
				}
			}

		case "reset":
			guild, _ := session.Guild(message.GuildID)
			if libs.CheckUser(session, guild, message) == true {
				var target models.User
				dbp.First(&target, "discord_id = ?", message.Mentions[0].ID)

				target.Balance = 0.00
				dbp.Save(&target)

				session.ChannelMessageSendEmbed(channel.ID, libs.ResetEmbed(message.Mentions[0].Username))
			}
		default:
			if len(message.Mentions) > 0 {
				getUser(message.Mentions[0].ID, message.Mentions[0].Username)
			}
		}
	} else {
		var person models.User
		if err := db.First(&person, "discord_id = ?", message.Author.ID).Related(&person.Houses, "Houses").Error; err == nil {
			session.ChannelMessageSendEmbed(channel.ID, libs.BankUserEmbed(message.Author.Username, person.Balance, len(person.Houses), false))
		} else {
			entry := models.User{
				DiscordID: message.Author.ID,
				Balance:   100.00,
				LastDaily: time.Now().Format(time.RFC3339),
			}

			if err := db.Create(&entry).Error; err == nil {
				session.ChannelMessageSendEmbed(channel.ID, libs.BankUserEmbed(message.Author.Username, 100.00, 0, true))
			} else {
				libs.Error(err, session, channel)
			}
		}
	}
}

func getUser(id string, name string) {
	var person models.User
	if err := dbp.First(&person, "discord_id = ?", id).Related(&person.Houses, "Houses").Error; err == nil {
		session.ChannelMessageSendEmbed(channel.ID, libs.BankUserEmbed(name, person.Balance, len(person.Houses), false))
	} else {
		libs.Error(err, session, channel)
	}
}

func checkUserBalance(input float64) bool {
	var person models.User
	dbp.First(&person, "discord_id = ?", "message.Author.ID")

	if person.Balance < input {
		return false
	}

	return true
}

func checkInput(input float64) bool {
	if input >= 0 || checkUserBalance(input) == true {
		return true
	}

	return false
}
