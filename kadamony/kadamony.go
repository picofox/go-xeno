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

	db.GetPoolManager().GetPool("DBP0").GetConnection(0).BeginTransaction()

	tlv1, ret := db.GetPoolManager().GetPool("DBP0").GetConnection(0).RetrieveField(db.DBF_TYPE_VARCHAR, true, false, "select token from account where uid=1")
	fmt.Println(ret)
	fmt.Println(tlv1)

	rd := db.NeoRecordDesc(0, false)
	rd.AddFieldDesc("uid", 11, db.DBF_TYPE_TIME, false, true)
	rd.AddFieldDesc("time", 11, db.DBF_TYPE_TIME, false, true)
	rd.AddFieldDesc("dt", 18, db.DBF_TYPE_DATE, false, true)
	rd.AddFieldDesc("ts", 18, db.DBF_TYPE_TIMESTAMP, false, false)
	rd.AddFieldDesc("creation_ts", 18, db.DBF_TYPE_DATETIME, false, true)
	rd.AddFieldDesc("nickname", 19, db.DBF_TYPE_VARCHAR, false, false)
	rd.AddFieldDesc("token", 20, db.DBF_TYPE_VARCHAR, false, true)
	tlv, rc := db.GetPoolManager().GetPool("DBP0").GetConnection(0).RetrieveRecord(rd, "select uid, time, dt, ts, creation_ts, nickname, token from account where uid = 2")
	fmt.Println(rc)
	fmt.Println(tlv.String())

	db.GetPoolManager().GetPool("DBP0").GetConnection(0).CommitTransaction()

	rData, _ := db.GetPoolManager().GetPool("DBP0").GetConnection(0).Retrieve("select dt from account limit 10")
	fmt.Printf("total %d lines", len(rData))
	for i := 0; i < len(rData); i++ {
		for j := uint32(0); j < rData[i].Length(); j++ {
			v, _ := rData[i].GetListValue(j)
			fmt.Print(v)
			fmt.Print(" - ")
		}

		fmt.Println("\n")
	}

	ra, rc := db.GetPoolManager().GetPool("DBP0").GetConnection(0).Delete("delete from account where uid = 10000")
	if rc != 0 {
		fmt.Println("delete failed")
	}
	fmt.Printf("delete %d \n", ra)
}
