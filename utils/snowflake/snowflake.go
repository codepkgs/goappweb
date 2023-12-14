package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
	"github.com/x-hezhang/gowebapp/settings"
)

var node *sf.Node

func Init(cfg *settings.SnowflakeConfig) (err error) {
	if err != nil {
		return
	}

	st, err := time.Parse("2006-01-02", cfg.StartTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(cfg.MachineId)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}
