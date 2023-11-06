package main

import (
	"fmt"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/db"
	"xeno/zohar/core/finalization"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/process"

	"xeno/zohar/core/datatype"
	_ "xeno/zohar/core/initialization"
)

func SetValues(f string, args ...any) {
	fmt.Println(f)
	fmt.Println(args)
}

func main() {
	defer finalization.GetGlobalFinalizer().Finalize()

	fmt.Println("Kadamony Intrinsic Initializing ... \t\t\t\t[Done]")
	kinfo := process.ProcessInfoString()
	fmt.Print(kinfo)

	netconfig := datatype.CreateTLV(datatype.DT_DICT, datatype.T_TLV, datatype.T_STR, nil)

	//netconfig.PathSet("  play[].nic[].ipv4.address ", "192.168.0.1", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)

	netconfig.PathSet("  default.nic[].ipv4.address ", "192.168.0.1", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[0].ipv4.address ", "192.168.0.2", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[].ipv4.address ", "192.168.0.3", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.gateway ", "192.168.0.1", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[1].id ", int32(998), datatype.DT_SINGLE, datatype.T_I32, datatype.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.dns[] ", "202.230.129.230", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.dns[0] ", "8.8.8.8", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.dns[] ", "111.11.11.111", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.metric ", 2, datatype.DT_SINGLE, datatype.T_I32, datatype.T_NULL)

	//_, rc := netconfig.PathGet("  default.nic[1].ipv4.metric ")

	fmt.Print("Kadamony Application Initializing ... ")
	rc, errString := config.LoadKadamonyConfig()
	if core.Err(rc) {
		fmt.Println(fmt.Sprintf("... \t\t\t\t[Failed] (%s)\n", errString))
	}

	db.GetPoolManager().Initialize(&config.GetKadamonyConfig().DB)

	fmt.Print("\t\t\t\t[Done]\n")

	rc = db.GetPoolManager().ConnectDatabase()
	if core.Err(rc) {
		logging.Log(core.LL_INFO, "Connect to Databases \t\t\t\t[Failed]")
	}
	logging.Log(core.LL_INFO, "Connect to Databases \t\t\t\t[Success]")

	table := db.NeoSQLTable("gamedb", "car", 0, true)
	table.AddFieldDesc("carId", 0, db.DBF_TYPE_BIGINT, false, false)
	table.AddFieldDesc("uid", 1, db.DBF_TYPE_BIGINT, false, false)
	table.AddFieldDesc("nftId", 2, db.DBF_TYPE_BLOB, false, true)
	table.AddFieldDesc("type_id", 3, db.DBF_TYPE_SMALLINT, false, false)
	table.AddFieldDesc("level", 4, db.DBF_TYPE_TINYINT, false, false)
	table.AddFieldDesc("part_level", 5, db.DBF_TYPE_BIGINT, false, false)

	var idx int64 = 0
	for idx = 0; idx < 1000; idx++ {
		table.AddRow([]any{100 + idx, idx + 1, nil, idx % 3, 1, 1122334455667788})
	}

	ps, rc := db.GetPoolManager().GetPool("DBP0").GetConnection(0).Create(table, "insert into car values(?,?,?,?,?,?)")
	fmt.Println(ps)
	fmt.Println(rc)
}
