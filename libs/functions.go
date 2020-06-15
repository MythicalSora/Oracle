package libs

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Error - This function sends an embedded error message if the command fails.
func Error(e error, session *discordgo.Session, channel *discordgo.Channel) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xff0000,
		Description: "[Error] => " + e.Error(),
		Title:       "Task Failed!",
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	session.ChannelMessageSendEmbed(channel.ID, embed)
}

// BankUserEmbed - This function returns an embedded message containing the user's bank information
func BankUserEmbed(user string, bal float64, houses int, new bool) *discordgo.MessageEmbed {
	var Title string
	if new == true {
		Title = "User created!"
	} else {
		Title = "Showing user " + user
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x90ee90,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Balance",
				Value:  fmt.Sprintf("%.2f", bal),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Houses",
				Value:  strconv.Itoa(houses),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     Title,
	}

	return embed
}

// BankPayEmbed - This functions returns an embedded message containing payment info
func BankPayEmbed(target string, bal float64) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x228b22,
		Description: ":money_with_wings: You paid $" + fmt.Sprintf("%.2f", bal),
		Title:       "Succesfully paid " + target + "!",
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	return embed
}

// ResetEmbed ...
func ResetEmbed(target string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0xfdee00,
		Title:     target + "'s balance has been reset.",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return embed
}

// AddEmbed ...
func AddEmbed(target string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0xfdee00,
		Title:     target + "'s balance has been increased.",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return embed
}

// SubEmbed ...
func SubEmbed(target string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0xfdee00,
		Title:     target + "'s balance has been lowered.",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return embed
}

// CheckUser ...
func CheckUser(session *discordgo.Session, guild *discordgo.Guild, message *discordgo.MessageCreate) bool {
	var roles []string
	var manager *discordgo.Role
	var owner *discordgo.Role

	for i := 0; i < guild.MemberCount; i++ {
		if guild.Members[i].User.ID == message.Author.ID {
			roles = guild.Members[i].Roles
		}
	}

	for i := 0; i < len(guild.Roles); i++ {
		if guild.Roles[i].Name == "Managers" {
			manager = guild.Roles[i]
		}

		if guild.Roles[i].Name == "Owner" {
			owner = guild.Roles[i]
		}
	}

	for i := 0; i < len(roles); i++ {
		if roles[i] == manager.ID {
			return true
		}

		if roles[i] == owner.ID {
			return true
		}
	}

	e := errors.New("you need to be Manager or higher to use this command")

	channel, _ := session.State.Channel(message.ChannelID)
	Error(e, session, channel)
	return false
}

// HousingBuyEmbed ...
func HousingBuyEmbed(user string, channel *discordgo.Channel) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x48c774,
		Title:       "House bought!",
		Description: "Congratulations " + user + "! Enjoy your new house, located at " + channel.Mention(),
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	return embed
}
