// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-12, by liasica

package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/liasica/edocseal/pb"
)

var w = sync.WaitGroup{}

func main() {
	num := 1
	w.Add(num)
	for i := 0; i < num; i++ {
		go do()
	}
	w.Wait()
}

func connect() (*grpc.ClientConn, pb.ContractClient) {
	conn, err := grpc.Dial(":17001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}

	// 建立gRPC连接
	return conn, pb.NewContractClient(conn)
}

func create() string {
	conn, client := connect()
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	jsStr := `{"address":{"Value":{"Text":"河南省济源市大峪镇偏看村"}},"aurDate":{"Value":{"Text":"2024年04月17日"}},"city":{"Value":{"Text":"成都市"}},"ebikeBattery":{"Value":{"Text":"时光驹电池"}},"ebikeBrand":{"Value":{"Text":"巨龙"}},"ebikeColor":{"Value":{"Text":"橘黄"}},"ebikeModel":{"Value":{"Text":"60V30AH"}},"ebikeSN":{"Value":{"Text":"2554219022000000"}},"ebikeScheme1":{"Value":{"Checkbox":true}},"ebikeScheme1PayMonth":{"Value":{"Text":"0"}},"ebikeScheme1PayTotal":{"Value":{"Text":"80.00"}},"ebikeScheme1Price":{"Value":{"Text":"80.00"}},"ebikeScheme1Start":{"Value":{"Text":"2024年04月17日"}},"ebikeScheme1Stop":{"Value":{"Text":"2024年04月26日"}},"idcard":{"Value":{"Text":"410881199504096034"}},"name":{"Value":{"Text":"王亚飞"}},"payMonth":{"Value":{"Text":"0"}},"phone":{"Value":{"Text":"17513170523"}},"riderContact":{"Value":{"Text":"[其他]疯狂 - 17566668888"}},"riderDate":{"Value":{"Text":"2024年04月17日"}},"riderSign":{"Value":{"Text":"王亚飞"}},"schemaEbike":{"Value":{"Checkbox":true}},"sn":{"Value":{"Text":"20240417114207523"}}}`

	var data map[string]any
	err := jsoniter.Unmarshal([]byte(jsStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	values := make(map[string]*pb.ContractFromField)
	for k, v := range data {
		value := v.(map[string]any)["Value"].(map[string]any)
		if check, ok := value["Checkbox"]; ok {
			values[k] = &pb.ContractFromField{
				Value: &pb.ContractFromField_Checkbox{Checkbox: check.(bool)},
			}
		} else {
			values[k] = &pb.ContractFromField{Value: &pb.ContractFromField_Text{Text: value["Text"].(string)}}
		}
	}

	res, err := client.Create(context.Background(), &pb.ContractCreateRequest{
		TemplateId: "C1455B9383BF453082B7341226EC60B3",
		Values:     values,
		Expire:     1713327127,
		Idcard:     "410881199504096034",
	})
	if err != nil {
		log.Fatalf("Call Create err: %v", err)
	}

	return res.DocId
}

func sign(docId string) {
	conn, client := connect()
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	res, err := client.Sign(context.Background(), &pb.ContractSignRequest{
		DocId:    docId,
		Image:    "iVBORw0KGgoAAAANSUhEUgAAAI4AAAA0CAYAAABcrAAbAAAABGdBTUEAALGPC/xhBQAAADhlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAAqACAAQAAAABAAAAjqADAAQAAAABAAAANAAAAADPJXAPAAAMo0lEQVR4Ae3ZCZAU1RnA8XAjN3Ir9y0QgYBQEMDlDARBQMQEkiBEikhColFjkRBiykiIlSgxRqiyTNAymNIcKhAIAUQUEBC5QRFhuRZYlPu+zP+/Tqe2Nuzs7DKzOwdf1Y+enenp7vf6e997PZT40vUorB4oy4l6ohH24HNcj+s9ELYHSvBpBxzEUlRCQkfJhL76xLn48lzqMNTEMZRBQkfxhL76xLj4YlzmjeiLs7Da1EVCR1Ekjudsjim4F3aipTxZw6reBBWwECZSJ1yPfPaAZXoqXBzqOKahMpIxbqBRo7AJw7EB/4LvxzJM2KYYi+/gNlRBwkY5rnw2TuD+0OuTbPshGcP2jsdm1MOvkIkRiGXU5uBP4hJcV+3DCjyHIaiGAkdRTFVerFPTRczDUziFZA372JFeGlbXF7Afd8CqEKvw8b8mVqMdPN9r8Fp+ilcwDq6/8h2xvPDcLsYG1UIGLmAylmIlkjFc0zg9224dxDOYhJawEsUiHJhWdSvOUaRjPbwelwXt0QdWpb9jPiKOwl6UetENMQEmThdUx0RkIhmjFI1qg65YhN2w8nSGC+ZliEV43sZojeU4hCDO8SIdq3Ae96AtdsBlQ55R2FOV57Pa1MdXYSKNhp2ZrOGIt3223eTxgeAAZmMQ6iIWcZaD7oLnb5TLCU7zvsuFX8Lfmv6AXsgzclYcG1cs9K1gm+dBctnBY8npUOXgXPtjNMaD8IKPIJnjCo3z5g1EfzhgBsBHdEe5sTTr3+j+4zmdInvAe7kEuYVTmdXH/b4PB/c6eIyrhjt2QhpuQV3YUNceljBL6mEcgjfY98/Azy6jNLw4f9SqDBdeZm6QMJZLk8b3XMW7v3PrVgxGsicNTcwK290TD6EBnCp8IGgN1zzDYJ9EO7yfLgNuxg+RV397L7vB69yIX8Mc+L8wcfphOBrBRnmjfd8y9hlce9g4s69sSFCpTDI7wX09gewQtyZYwLLpsfzeWvhkMRkmXyqF/erNsS9bYgrsUwfiWJxENMMBOxRjYAJFkpzeIxPaRPP1E9iBq4YNsmKYmXVRD3VQFTYsewTVJJiC/G4k4f4DYeeMiOQLSbqP/dAej6ETOsDR/VuYVNEMz9Uf76BzPg7sPW2KmXgVrVBkYQJOgws2nzJSMSrS6HswCbeGOsBkGYItmAEHcLTCitETC9ClAAetz3emwwW0yV4kYdlcjKVw7ZMqYYWuBRPG6XkQfFDIHlYGn7gWYg4c7dEIE6cb/omCJI7X0ASz4DFcymSFjSqMsPQ5/dmBZv9FJHs4OJyGnJZH4TCehYlxBtnjEn+sxHfxKZwifAK91vicA5zDZ4h0SZHznDt5w0pYCaNhMn7xjy9iHI6ozrgLlmkbkqzhlNwbVhAH5lbMhVO0NzFcnOBDq7JT2gPYD29cQcPEcepriU+wGwWJk6Ev+fS3CfuzsqcgR8rnd0ycO1AFjrorSKawfT5YuNh1asjAWnyA/A4Sf/LYgPOYiOrYgZxVirciChPHpySvyeQpSHgtx9Ecrs0WFVbiVOBk93lCWJKTIawmNdAWjugW8Ob8G47sYJTyMt/h1LUN29EXw2El24uzyE84WH0YOQATsKBhtbTNPbDGkZJXODf6BcteEPmpGH6/GhyRS5GoYTu8ec71rj9uhNXE37m8oYpmuA5cjtXoiW9gKFykzsFRRBLe49K4HMnOYfaxEh6ECVQjXOKYLI6mB2HFCDLOBNqHJTCDHVmOAhvqPjmTyqrmr9IuDp3vEylMlqqoDxPGacO2+rvLMdhmOzSW4TkWYgV64ZsYgL/hLRxBbuE9rAiT3AS/lggGjjlzIVziuKNlrgMsyT5C3gAvxs56CF6MTB7fs1NNnlOwY22U235wodcEfnYaJ5C9ivFn3ISVpR2+DK/fpLeiLENRhf32Jhyw3XEXRsGk8v0M5AzvWxr87oe4ljD5XMN5n7eaHOHCzz15GVg5ZHjDa2IcRsLOfRzpqI7sSVaXvy2xs2FpLw8rmMeUx3JU2TgXgF6Yi7Hgb5PRKmaS+12/Y/L6Pa/Hrd/xuzIh/a7Jafi5YamWx/L4ntMk9zuGg6I2OqIRPsJGeM3HEW9hJemCwbgZ67AYu2HYLiv9BEzHGhQ07JcfIA3TMDevxGGfsOHN7I/fwBtkI3YiiFK8GAET524EN5GX//tdwSQwIRzlJoVbeSPdlobnsSNMBpPEG+6xvH4/d1870lFRKbQ1Of3c5HJrwhgeRy5Afa8eTKQmodcz2L6GRAjb3gL27xjYFqu//bMNe2A7f4ZDyE/YbzVxKwahIZ7BQmR1qNtrCW98a9jZq3E/gtFu5XkWm/A04iFMRqvgQHTGLrwHq1RlOIoz8BKOIt7CQeAN7Q8HqtXAZcB6bIGJYwJZaRyIb8OB5SC2eh4IbX1tMmXCtntc298MJksbOPjsi7WYD7+bFe4cjbAiOOc+hwEwgTx2A8zEIzB5iipMbsu5HXILqmMlliF7cnjNtTAc7WGHvYLs+/BnkURwbSZMT3hN87AOTs9WYsPK+Shs85NIRxmYPLb7JlRFndBrk9AEMaxOJuFH8H5th0UgODYvvwhveDTiMgdZhT3oDhPHUudN8kKyT1/8WWhRnjN5PV3hdJSOBdiMq4XT30G8DKtSb1hJizqcgu9EX3wCpwyTOnuYWK3wExzBVByG4Y23EmViK645opU4dvgpOAqah66qFFsb6jRg1hZWVOFEXkNHmLgf4w2kw/KcVzj6RqIcxuNTFFWYDCawyeBAfAobkbMCWF1Mmkk4jicQ0+uOVuJwnVkj+ixbpwQbXA025nnEOuy41nB6saMty+/iRUSStH7fKaod0lARjyGmnc/xw4UDrzPGYQ1mwcGZM8ryRje4n+uRaYj5dUczcXxCccqqAI/bB47wDxGLsGNdGPZAF+zDBrwDF7zhwtFbCY3RCx1ghdmO97EIMe98zpFb2IdjkIY/Yy5yhn3cFO5ndf0HXsA5xDw8ebTC6eoiLqE0+sIG+360wkrWEN5oR6NVbQUexyEE4X5WDdcGVUPc16ri1vdugtfm4nIGfCqxOnn9RRVet4vViaiDh7EL2cMB0wzD4aDZDPePytqF40QU+U0cG+YawtFqp8v3zoZeW3V83QTeHEeuURx+ljP8bm7hsYOwHN+OoagMq9g8rIWj0/OloQEawk43ec/AZAi2mbzeBUu/x9iJeAn7yAoyBbvxAE4iCPu9NwbDvjXRTSy3hR7hbtzVLqYGb76IjjARgu+f57XTlCPchZuLZBuaDqcAO+UIrA4HYcL6ntuAfxsex2N44+uhDDyWnVoC3nATwGSyA00wz+d7++Hx03EAJkzAa4zXsIp0gZVjCWbCtrVBP7SC/XgcDpi3sRfegyKJ4MZHenJv5lB4Ey3pNtgb7c3xBt4LpwOnrD/hPZhM5dELfncCPsA5eDMDQcJYzfpjCEwcjzEXh2F4bPe1khyDnZnI4cCwb36E/8D2dofvWV03YRm2wWnJip5U4Yh4EzZsC0yuIEyw8TiBZsGb2bYmcCM44jzGXzACdZDM4dPfWFgx07EeDgYTZCRchyV9WCkWwIrwOhxJQfjZX7EaJlgQVqM++B0WYjq6wsqW7GE/jMYqWIFd3CsTJlIG/ogWiLsoGcUrcv3hDXcKexdOQUZxNEYaJuMCBuNONMdBzMfT2IdUCCvNMHwL38N+WF3tG/vQAdUO96Euvo0TSMrwMXc5fBK4LVsLrTy/gEnxMFzY7YAVqCesRqkUJo2JYD9YXXML+8014R50yW2nZHjfkTIeU3I0pjJ/W4odTY6aVxGuw/g4acNk+DqWYhDChZW6BXzKGhtux2T7rBgNagIriwvmxUiDU1oqhkljsvjk1DeCDjBxnK78zcnKkxJho9vCpyuTxjm8H1I1nJ5G4n24voskTLQhcIr395yUCNctc2DCOHWZQAORilGKRlth/OHShXCk4ZrxJVihnO6TPqw2Lo4P4HnUxjb8HKkWTtf18Qbsi0inaZPta7DaPIKUCBOnNXZiPR6FHeDI8beLVArbOwFWG/skkjC5HHg+ULwOH8dTJmx8f6zFFfjfAjPhk1cqRUMa6+CZCqtPXlGeHcbjY7yMWkjJsCO+EuL/YqdSlKWxd8PKkdfvMCXZpx1mYQ/83ct1YtyGFxzLOM3B7bhUDKfsijgHK+7Vws/9TWs4boc/jA6DT19xHbFOnLhufIwvzh88TQSfiPxpIjP0ujFbdYLrHqcw93Mt6DrwFOI+Ipl3474RcXyBPlH+HlaV3XCxbJ9nYCvegtXFpEqouJ44sb9dTkfNYdL4Xy57Q9tLbBM2/guFvdKfumi2aQAAAABJRU5ErkJggg==",
		Name:     "王亚飞",
		Province: "成都市",
		City:     "成都市",
		Address:  "河南省济源市大峪镇偏看村",
		Phone:    "18888888888",
		Idcard:   "410881199504096034",
	})
	if err != nil {
		log.Fatalf("Call Sign err: %v", err)
	}
	fmt.Println(res)
	w.Done()
}

func do() {
	docId := create()
	sign(docId)
}
