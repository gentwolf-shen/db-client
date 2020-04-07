package DbClient

type (
	SqlMessage struct {
		Token    string        `json:"token"`
		Database string        `json:"database"`
		Sql      string        `json:"sql"`
		Params   []interface{} `json:"params"`
	}

	BatchSqlMessage struct {
		Token    string          `json:"token"`
		Database string          `json:"database"`
		Sql      string          `json:"sql"`
		Params   [][]interface{} `json:"params"`
	}
)
