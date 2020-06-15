package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mythicalsora/Oracle/commands"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mythicalsora/Oracle/models"
)

var db *gorm.DB
var dbe error

var user models.User = models.User{}

// Config ...
type Config struct {
	Token      string `json:"token"`
	Prefix     string `json:"prefix"`
	ConnString string `json:"connString"`
}

var config Config

func createTables() {
	db.AutoMigrate(&models.User{})
}

func init() {
	db, dbe = gorm.Open("mysql", config.ConnString)
	if dbe != nil {
		fmt.Println("Couldn't conenct to Database:", dbe)
	}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("There was an error reading your config file, ", err)
		return
	}

	err = json.Unmarshal([]byte(string(data)), &config)
	if err != nil {
		fmt.Println("There was a problem decoding your config file, ", err)
	}

	createTables()
}

func main() {
	bot, err := discordgo.New("Bot " + config.Token)
	check(err)

	bot.AddHandler(onMessage)
	err = bot.Open()
	check(err)

	fmt.Println("Client is up & running! Press CTRL-C to quit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("\nCleaning up...")
	bot.Close()

	defer db.Close()
}

func check(e error) {
	if e != nil {
		fmt.Println("There was a problem: ", e)
	}
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID || message.Author.Bot {
		return
	}

	if strings.HasPrefix(message.Content, config.Prefix) {
		channel, err := session.State.Channel(message.ChannelID)
		msg := strings.TrimPrefix(message.Content, config.Prefix)

		args := strings.Split(msg, " ")
		command := args[0]
		args = append(args[:0], args[1:]...)

		check(err)

		switch command {
		case "ping":
			fmt.Println("test:", channel)
		case "bank":
			commands.Bank(session, message, db, channel, args)
		case "housing":
			commands.Housing(session, message, db, channel, args)
		default:
			fmt.Println("Command: ", command)
			fmt.Println("Args: ", args)
			fmt.Println("Full message: ", msg)
		}
	}
}
