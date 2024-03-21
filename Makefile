.PHONY:  start run clean killport	kill	start
#后台启动,持续工作
start:
	nohup ./collection > collection.log 2>&1 &
#开启容器，然后编译程序
run:
	cd fixtures && docker-compose up -d && cd .. && go build -o collection
#删除容器和编译文件
clean:
	cd fixtures && docker-compose down -v && cd .. && rm -f collection
#杀死占用8080的进程
killport:
	kill -9 $$(lsof -t -i:8080)
#杀死程序进程
kill:
	pkill -9 collection

