package main

import (
	"fmt"

	"github.com/datasweet/datatable"
)

func main() {
	dt := datatable.New("test")
	dt.AddColumn("champ", datatable.String, datatable.Values("Malzahar", "Xerath", "Teemo"))
	dt.AddColumn("champion", datatable.String, datatable.Expr("upper(`champ`)"))
	dt.AddColumn("win", datatable.Int, datatable.Values(10, 20, 666))
	dt.AddColumn("loose", datatable.Int, datatable.Values(6, 5, 666))
	dt.AddColumn("winRate", datatable.Float64, datatable.Expr("`win` * 100 / (`win` + `loose`)"))
	dt.AddColumn("winRate %", datatable.String, datatable.Expr(" `winRate` ~ \" %\""))
	dt.AddColumn("sum", datatable.Float64, datatable.Expr("sum(`win`)"))

	fmt.Println(dt)
}

/*
CHAMP <NULLSTRING>      CHAMPION <NULLSTRING>   WIN <NULLINT>   LOOSE <NULLINT> WINRATE <NULLFLOAT64>   WINRATE % <NULLSTRING>  SUM <NULLFLOAT64>
Malzahar                MALZAHAR                10              6               62.5                    62.5 %                  696
Xerath                  XERATH                  20              5               80                      80 %                    696
Teemo                   TEEMO                   666             666             50                      50 %                    696
*/
