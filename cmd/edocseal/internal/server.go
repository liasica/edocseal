// Copyright (C) edocseal. 2024-present.
//
// Created at 2024-04-07, by liasica

package internal

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/liasica/edocseal/internal/g"
	"github.com/liasica/edocseal/internal/service"
	"github.com/liasica/edocseal/internal/task"
	"github.com/liasica/edocseal/pb"
)

func serverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "server",
		Short:             "启动服务端",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Run: func(_ *cobra.Command, _ []string) {
			// 启动任务队列
			task.CreateTasks(g.GetSignTaskNum(), g.GetDocumentTaskNum())

			// 启动http服务
			go func() {
				http.HandleFunc("/stop", service.StopTasks)
				zap.L().Info("API启动", zap.String("bind", g.GetHttpBind()))
				err := http.ListenAndServe(g.GetHttpBind(), nil)
				if err != nil {
					fmt.Printf("HTTP服务启动失败：%s\n", err)
					os.Exit(1)
				}
			}()

			// 监听端口
			lis, err := net.Listen("tcp", g.GetRPCBind())
			if err != nil {
				fmt.Printf("监听TCP端口失败：%s\n", err)
				os.Exit(1)
			}

			// 创建grpc server
			s := grpc.NewServer(
				grpc.ChainUnaryInterceptor(
					// logging.UnaryServerInterceptor(edocseal.InterceptorLogger(zap.L())),
					recovery.UnaryServerInterceptor(),
				),
			)
			defer s.GracefulStop()

			pb.RegisterContractServer(s, &service.ContractService{})
			zap.L().Info("RPC启动", zap.String("bind", g.GetRPCBind()))

			// 启动服务
			err = s.Serve(lis)
			if err != nil {
				fmt.Printf("RPC服务启动失败：%s\n", err)
				os.Exit(1)
			}

			select {}
		},
	}

	return cmd
}
