package todo

import (
	"github.com/boltdb/bolt"
	"github.com/gopherworks/bawt"
)

type Plugin struct {
	bot   *bawt.Bot
	store Store
}

func init() {
	bawt.RegisterPlugin(&Plugin{})
}

func (p *Plugin) InitPlugin(bot *bawt.Bot) {
	p.bot = bot

	err := bot.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
	if err != nil {
		p.bot.Logging.Logger.Fatalln("Couldn't create the `todos` bucket")
	}

	p.store = &boltStore{
		db:  bot.DB,
		log: bot.Logging.Logger,
	}

	p.listenTodo()
}
