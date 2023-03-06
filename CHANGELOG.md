# CHANGELOG

## [v2.3.1](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v2.3.0...v2.3.1)

> 2023-03-06

### Bug Fixes (1)

* **trade:** fix last tick always nil, rename eventbus subscribe naming

### Code Refactoring (1)

* **common:** rename common to global

## [v2.3.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v2.2.0...v2.3.0)

> 2023-03-06

### Bug Fixes (5)

* **ci:** fix make build fail in go test
* **dt:** remove redundant to check tick rate
* **grpc:** extend connection timeout to 30 second in every try connect
* **health:** make sure both sinopac and fugle are all down then terminate
* **trade:** fix wrong useage of in-out-ratio, wrong max hold time in simulate

### Features (4)

* **cron:** add publish terminate from rabbitmq when health check fail or in cron job
* **grpc:** move development to config, add client id, add beat error message
* **simulation:** add simulation to history router
* **simulation:** add simulation for future back, extend max idle time for postgres

## [v2.2.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v2.1.0...v2.2.0)

> 2023-02-24

### Bug Fixes (4)

* **ci:** fix missing env file in build test
* **ci:** fix test missing env file
* **order:** fix cancel order always show cancel fail
* **slack:** remove redundant in cancel process message

### Code Refactoring (1)

* **trade:** redesign trade rate usage

### Features (3)

* **slack:** add slack module instead of implement to grpcapi
* **slack:** add place and filled message for slack
* **trade:** add back check current order is cancelled or not

## [v2.1.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v2.0.1...v2.1.0)

> 2023-02-23

### Bug Fixes (9)

* **badge:** fix wrong badge link
* **event:** fix wrong event time due to wrong database time zone
* **order:** fix cancel order will continue if returen is id already cancelled
* **slack:** add InsecureSkipVerify to fix verify fail
* **slack:** fix wrong place of notify place order
* **trade:** fix same pointer cause place redundant order
* **trade:** fix cancel order multi times
* **trade:** fix wrong use of cancel order time
* **trade:** fix wrong trade calculate method

### Features (11)

* **basic:** add option basic data
* **config:** move rate_limit, rate_change_ratio to config
* **log:** remove slack hook for logger, modify cancel order method
* **order:** move modify future night market order time to sinopac
* **pkg:** change pkg name from topic to event
* **slack:** add logger hook to slack
* **slack:** add notify to slack when buy sell cancell future order
* **target:** filter stock target if it has its own future
* **trade:** add cool down time 3 minutes between trade in
* **trade:** add check gap between tick if is lower than 1 second
* **trade:** remove cancel wait time, use buy sell wait time instead

## [v2.0.1](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v2.0.0...v2.0.1)

> 2023-02-09

### Bug Fixes (1)

* **trade:** fix check wait time fail when last tick is nil

## [v2.0.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.7.0...v2.0.0)

> 2023-02-09

### Bug Fixes (22)

* **assist:** fix not trade out
* **assist:** fix does not try get more balance
* **basic:** fix dupl pkey in update future basic data
* **ci:** fix golang version missing double quote
* **ci:** fix wrong golangci-lint version
* **event:** fix wrong sinopac event time, change package naming mq to rabbit
* **events:** fix stock trade room will start before history data fetch done
* **future:** fix ws client get no future tick
* **healthcheck:** fix wrong usage of recover
* **index:** fix yahoo price is nil
* **index:** fix if snapshot or yahoo price is nil cause panic
* **index:** if nasdaq, nf is zero, not return error but log warning
* **index:** cancel multithreading in get index, yahoo price has new error
* **lint:** fix ci lint version
* **order:** fix order time always now cause will not be cancelled in night market
* **realtime:** fix stuck when new hadger
* **target:** fix stock future target mixed, refactor target topic
* **trade:** fix mem leak in check balance, fix order time in rabbit, add cancel order in day trader
* **trade:** fix balance trader missing equal symbol
* **trade:** fix wrong trade time in ws future trade
* **trade:** fix alway get fixed balance
* **wording:** fix wrong wording in last commit

### Code Refactoring (5)

* **cache:** refactor cache pkg and module, rename modules to module
* **db:** modify initail db method
* **event:** remove event module, refactor event bus pkg
* **grpc:** refactor attampts method, rename events to topic
* **usecase:** re-design rpc, add realtime usecase, let method split make sense

### Features (36)

* **assist:** modify trade out method to increase more profit possible
* **assist:** add trade out price to get more balance
* **balance:** modify judge forward or reverse order method
* **config:** redesign config naming
* **dt:** finish dt module without generate order method
* **fugle:** add fugle grpc to order and basic alpha
* **grpc:** split sinopac and fugle grpc api
* **hadger:** add rabbit for hadger
* **healthcheck:** recover panic if manual stop toc-sinopac-python, update go
* **history:** remove biasrate, update to latest proto, realtime usecase use own rabbit
* **index:** improve get all index method performance
* **interface:** make sure all new instance return interface instead of original object
* **logger:** replace all panic to fetal
* **order:** remove tick time in stock, future order
* **order:** decrease future trade cancel time to 5 seconds
* **order:** remove trade time in order
* **order:** remove manual, group id, add dt trader beta
* **pkg:** refactor log, config package
* **position:** add limit get position in ws stream, if not trade time, return
* **proto:** use new format of proto
* **rabbit:** let one trader has own rabbit connection
* **router:** add get trade index api
* **router:** add last stock or future balance router
* **simulate:** modify simulate method, fix helath check router
* **trade:** add hold times and trade out price
* **trade:** increase check times, fix heartbeat panic fn
* **trade:** finish first version dt for future
* **trade:** move try get more balance out of high or low
* **trade:** add support more than 1 qty in assist trade
* **trade:** add max hold time, and check switch in dt
* **trade:** add notify trade switch in dt and hadger
* **trade:** add trade out wait times in dt
* **trade:** add balance high, low to day trade future, finish future day trader unit
* **trade:** remove switch router, add health check router, change trade index to index status
* **trader:** add day trader for future alpha
* **usecase:** add usecase base, add hadger alpha, split interfaces, config read once

## [v1.7.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.5.0...v1.7.0)

> 2023-02-09

### Bug Fixes (33)

* **assist:** add lock for process trade, fix bugs in assist trader
* **assist:** fix assist trader will kill before start
* **ci:** add checkout in deployment to get git hash
* **ci:** add missing config in go test
* **ci:** temp remove go test from ci
* **future:** add missing future detail
* **kbar:** fix wrong time period of stream kbar, add send kbar every minute
* **lint:** fix stream routes lint error
* **log:** remove unknown order code log
* **order:** modify get order by trade day will send not filled order
* **order:** fix order status chan is not add to rabbit, fix order time is always wrong
* **order:** fix order balance calculate wrong, try fix manual order does not insert to db
* **order:** fix manual order does not insert to db
* **order:** add lock for order usecase to update order in postgres
* **order:** modify send order sequence to avoid stuck
* **order:** use timer and reset to fix balance not insert to db
* **position:** fix websocket future position has no column in json
* **postgres:** fix redundant manual future trade balance cause panic
* **router:** fix return body of get future trade switch
* **stream:** add loop lable to avoid index out of range
* **stream:** fix last trade rate not initial
* **stream:** fix map is not initail
* **stuck:** fix missing go in updateAllTradeBalance
* **tick:** fix wrong tick time from snapshot to tick
* **trade:** fix order will be cancel multiple times
* **trade:** fix order will be cancel before 10 seconds
* **trade:** fix wrong first tick time
* **websocket:** fix order map does not initial
* **websocket:** fix send data to close channel
* **ws:** fix stuck if client disconnect by abort gin
* **ws:** remove lock for future stream
* **ws:** fix wrong ws message type
* **ws:** fix ws end abnormal, close connection before gin done

### Code Refactoring (4)

* **assist:** refactor assist trader
* **logger:** pack log again
* **ws:** split pick stock and future to different pkg
* **ws:** refactor websocket split pick stock and future trade

### Features (42)

* **assist:** add assist trader in futrue trade ws
* **assist:** let buy sell has same profit loss automation method
* **assist:** add assist done message, limit 1 assist at one time
* **assist:** send running status to ws
* **assist:** finish by balance and by time period assist
* **balance:** add move order and recalculate trade balance router
* **balance:** add manual balance router
* **balance:** add manual trade balance
* **future:** split kbar and send last period to stream
* **future:** add send future detail in first connect
* **index:** add otc, tse, nasdaq index to websocket stream
* **kbar:** add try last day if kbar is empty, slow down query position
* **log:** move log config to env instead of in config
* **order:** add if not stock and not future trade time, cancel update order
* **order:** refactor updateAllTradeBalance, increse update order speed
* **order:** remove non block order update mode
* **order:** split prod and simulate order get method
* **order:** add receive order arr from rabbitmq
* **order:** add if simulate mode to get order status from mq
* **order:** remove ask update api, split update balance and simulate order, product order
* **order:** modify new sinopac filling to partfilled
* **order:** remove order status in websocket, modify manual order group id and trade time
* **order:** add order status stream in websocket
* **ordertime:** modifiy wrong order time to time now, when in night market from 0:00 to 5:00
* **position:** increase speed of send future position to every second
* **postition:** add future postion, remove manual in stock, future balance
* **protobuf:** change future trade all json message to proto message
* **router:** add manual insert future order api
* **router:** add query future order by date api
* **snapshot:** add future detail in future snapshot
* **stream:** add out in rate chg, modify order entity
* **stream:** add period update trade index interal stream usecase
* **stream:** change method process tick arr
* **stream:** add kbar in future trade ws
* **stream:** modify out in volume to four period
* **subscribe:** add whether subscribe stock or future not
* **switch:** add change future trade switch router
* **tick:** cut tick arr in every second in stream
* **trade:** add manual to order column, check order status in websocket
* **trade:** add auto cancel in stream trade
* **websocket:** add nasdaq future to stream
* **ws:** add lock for future trade ws

## [v1.5.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.4.0...v1.5.0)

> 2023-02-09

### Bug Fixes (10)

* **health:** if disconnect from grpc not panic but os exit
* **migration:** fix wrong table name in migration sql
* **subscribe:** remove redundant future bidask subscribe
* **websocket:** remove period to fix out of index
* **websocket:** fix wrong scoket data type
* **websocket:** fix wrong calculation of trade rate
* **websocket:** fix missing format in send websocket data
* **websocket:** add missing socketPickStock
* **websocket:** fix interface cast bug
* **websocket:** fix concurrency write websocket

### Features (12)

* **action:** use period tick out in ratio to decide action, and add to simulator
* **history:** add future history close fetch and simulate
* **holiday:** extend trade year to 2023, update holiday.json
* **order:** modify all order router to future and stock all order
* **tick:** modify future tick chan and connection id
* **trade:** add ws trade
* **tradeperiod:** add get last 1 trade period for future method
* **websocket:** modify trade rate content to out, in and period
* **websocket:** add trade rate in websocket
* **websocket:** ignore CloseNoStatusReceived error
* **ws:** add log for new future ws and done log
* **ws:** modify ws layout, add send snapshot in future stream ws

## [v1.4.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.3.0...v1.4.0)

> 2023-02-09

### Features (8)

* **future:** use rsi to decide trade in, max hold time to trade out
* **future:** test new method of future trader
* **future:** modify trade out method
* **future:** add kbar analyze to trade out of future
* **router:** add simulate future to history router
* **simulate:** add simulate future trade proto type
* **strategy:** modify future trade strategy default config
* **trader:** add kbar analyze to future trader

## [v1.3.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.2.0...v1.3.0)

> 2023-02-09

### Features (3)

* **balance:** exclude unfinish order in calculate balance, modify topic name
* **router:** add future trade balance to get balance
* **strategy:** modify future trader strategy

## [v1.2.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.6.0...v1.2.0)

> 2023-02-09

## [v1.6.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.1.0...v1.6.0)

> 2023-02-09

### Bug Fixes (35)

* **assist:** fix assist trader will kill before start
* **assist:** add lock for process trade, fix bugs in assist trader
* **ci:** temp remove go test from ci
* **ci:** add checkout in deployment to get git hash
* **ci:** add missing config in go test
* **health:** if disconnect from grpc not panic but os exit
* **log:** remove unknown order code log
* **migration:** fix wrong table name in migration sql
* **order:** fix manual order does not insert to db
* **order:** modify get order by trade day will send not filled order
* **order:** use timer and reset to fix balance not insert to db
* **order:** add lock for order usecase to update order in postgres
* **order:** fix order status chan is not add to rabbit, fix order time is always wrong
* **order:** fix order balance calculate wrong, try fix manual order does not insert to db
* **position:** fix websocket future position has no column in json
* **postgres:** fix redundant manual future trade balance cause panic
* **router:** fix return body of get future trade switch
* **stream:** fix map is not initail
* **stream:** fix last trade rate not initial
* **stream:** add loop lable to avoid index out of range
* **stuck:** fix missing go in updateAllTradeBalance
* **subscribe:** remove redundant future bidask subscribe
* **trade:** fix order will be cancel multiple times
* **trade:** fix order will be cancel before 10 seconds
* **websocket:** remove period to fix out of index
* **websocket:** fix send data to close channel
* **websocket:** fix wrong scoket data type
* **websocket:** fix wrong calculation of trade rate
* **websocket:** fix missing format in send websocket data
* **websocket:** add missing socketPickStock
* **websocket:** fix interface cast bug
* **websocket:** fix concurrency write websocket
* **websocket:** fix order map does not initial
* **ws:** fix ws end abnormal, close connection before gin done
* **ws:** fix stuck if client disconnect by abort gin

### Code Refactoring (5)

* **assist:** refactor assist trader
* **logger:** pack log again
* **trader:** put stock future trader into a module, change the stream usecase
* **ws:** split pick stock and future to different pkg
* **ws:** refactor websocket split pick stock and future trade

### Features (51)

* **action:** use period tick out in ratio to decide action, and add to simulator
* **assist:** let buy sell has same profit loss automation method
* **assist:** add assist done message, limit 1 assist at one time
* **assist:** finish by balance and by time period assist
* **assist:** add assist trader in futrue trade ws
* **assist:** send running status to ws
* **balance:** add manual balance router
* **balance:** add move order and recalculate trade balance router
* **balance:** exclude unfinish order in calculate balance, modify topic name
* **balance:** add manual trade balance
* **future:** add kbar analyze to trade out of future
* **future:** use rsi to decide trade in, max hold time to trade out
* **future:** test new method of future trader
* **future:** modify trade out method
* **history:** add future history close fetch and simulate
* **holiday:** extend trade year to 2023, update holiday.json
* **index:** add otc, tse, nasdaq index to websocket stream
* **log:** move log config to env instead of in config
* **order:** modify all order router to future and stock all order
* **order:** modify new sinopac filling to partfilled
* **order:** add order status stream in websocket
* **order:** refactor updateAllTradeBalance, increse update order speed
* **order:** remove ask update api, split update balance and simulate order, product order
* **order:** remove order status in websocket, modify manual order group id and trade time
* **ordertime:** modifiy wrong order time to time now, when in night market from 0:00 to 5:00
* **postition:** add future postion, remove manual in stock, future balance
* **router:** add query future order by date api
* **router:** add manual insert future order api
* **router:** add future trade balance to get balance
* **router:** add simulate future to history router
* **simulate:** add simulate future trade proto type
* **strategy:** modify future trade strategy default config
* **strategy:** modify future trader strategy
* **stream:** add period update trade index interal stream usecase
* **stream:** add out in rate chg, modify order entity
* **stream:** change method process tick arr
* **stream:** modify out in volume to four period
* **switch:** add change future trade switch router
* **tick:** cut tick arr in every second in stream
* **tick:** modify future tick chan and connection id
* **trade:** add auto cancel in stream trade
* **trade:** add manual to order column, check order status in websocket
* **trade:** add ws trade
* **tradeperiod:** add get last 1 trade period for future method
* **trader:** add kbar analyze to future trader
* **websocket:** add trade rate in websocket
* **websocket:** ignore CloseNoStatusReceived error
* **websocket:** modify trade rate content to out, in and period
* **websocket:** add nasdaq future to stream
* **ws:** add log for new future ws and done log
* **ws:** modify ws layout, add send snapshot in future stream ws

## [v1.1.0](https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.0.0...v1.1.0)

> 2023-02-09

### Bug Fixes (1)

* **future:** fix wrong future trade switch, add future trade fee calculate

### Features (5)

* **future:** add auto select r1 mxf future
* **module:** move event topic to a new event module with bus
* **module:** move target filter to a new target module
* **module:** move all trader to trader module, remove simulate stock trader
* **tradeday:** add new future trade day method, modify query order method to time range

## v1.0.0

> 2023-02-09

### Bug Fixes (32)

* **actions:** add reset environment before deployment
* **cancel:** fix cancel fail casue filled order not append to order map, add cancel wait time
* **ci:** fix wrong way to load env file
* **config:** fix simulate must be true
* **cpu:** fix checkFirstTickArrive cause cpu 100%
* **deployment:** fix wrong config path in deployment action
* **event:** fix missing subscribe bidask
* **future:** fix wrong rsi gap base, rename future trader and stock trader
* **history:** fix skip close is 0 in insert db panic
* **history:** fix wrong key with new fetch kbar tick
* **history:** no gorutine when process data to avoid unexpected error
* **logger:** fix wrong use format in log error
* **logger:** fix debug missing format function
* **open:** fix wrong open change ratio in trade
* **order:** update order in every order status return, no compare
* **order:** fix wrong out in ratio in order generator
* **order:** fix repeat place order and cancel
* **order:** fix empty order time when update, temp extend trade in, out wait time
* **order:** fix wrong status when order updated, add last tick in trader, add lock in cancel order
* **order:** using waiting order in tick for, remove uuid, place order fail add status
* **path:** modify initial sequence, fix wrong event subscribe callback
* **quota:** fix wrong quota when sell or buylater, fix check cancel order wrong tarde time
* **quota:** fix wrong calulate quota, check order status is failed, add lock in place order
* **readme:** fix wrong attachment path
* **simulate:** fix wrong orders return, modify default config, add realtime target to 30 secs
* **snapshot:** fix all snapshot return empty panic, insert all stock from sinopac
* **status:** fix fetch history done event wrong input function
* **table:** fix wrong table name fail to create in postgres
* **target:** fix wrong repo table name when update target
* **target:** fix stuck by non async event, add fetch list in history usecase
* **trade:** fix wrong analyze tick time, if first tick not arrive no action
* **tradeagent:** fix wrong in out ratio compare

### Code Refactoring (6)

* **config:** move config pkg to top level between cmd
* **dependency:** remove logger dependency in all pkg
* **logger:** remove global basepath set and get, remove global dependency from logger
* **pkg:** rename sinopac to grpc, move global to top level
* **pkg:** refactor config, eventbus method
* **target:** rename subscribe and realtime add to pre-fetch, realtime

### Features (104)

* **agent:** add new agent method, reanme trade to trade agent
* **analyze:** add last trade day all ticks cahce, upgrade go to 1.18.4
* **analyze:** add history open from kbar, fix wrong routing key of bidask, redesign cache
* **analyze:** add biasrate of history close
* **analyze:** add quater ma to cache
* **analyze:** add below quater ma stocks, include api, add history analyze table
* **analyze:** split stock, future analyze, fix future order repo, future first tick analyze time
* **analyze:** same way to analyze realtime tick,history tick, order trade time assign in place
* **api:** add day trade calculator, add config api
* **balance:** add all trade balance api
* **balance:** add calculate future order balance, add allow trade in stock, future config
* **basic:** skip add category is 00 to stock list
* **basic:** add import holiday from data/holiday.json
* **basic:** add insert or update stock, remove stock entity id, to number
* **biasrate:** change the way of biasRate usage, use bidask price to decide trade out price
* **cache:** add stock detail cache in basic usecase
* **cache:** refactor cache, add basic info, fix insert sql too many parameters
* **change-ratio:** consider open change ratio to unsubscribe, and change ratio in stream
* **config:** add tag for config api, rename most api method
* **config:** add sinopac path to config and env, modify readme env template
* **config:** remove redundant config, change default config trade fee discount ratio
* **config:** add trade about config, change terminate to put
* **entity:** add base order to distinguish stock order and future order
* **event:** add eventbus package, rename from stock to basic
* **event:** add log if event is not about subscribe
* **event:** remove fetch_history_done event
* **future:** use future gap in 8:45 to decide forward or reverse, add trade day struct
* **future:** add future detail, add subscribe future tick
* **go:** update go to 1.19.1 and remove k8s agent
* **grpc:** move snapshot from basic to stream, add tse snapshot api, add snapshot entity
* **heartbeat:** add heartbeat, history entity modify logger initial
* **history:** add day kbar api, add all usecase router
* **history:** add history tick analyze, and use for realtime order generator
* **history:** modify method to check already exist kbar, tick, merge stream tick, bidask processor
* **history:** add history close fetch, add terminate sionopac api, remove open in history close
* **history:** add delete before fetch
* **history:** add fetch history event lock
* **kbar:** add fetch history kbar
* **log:** merge all log into one file, and pretty json
* **logger:** modify log format and check if development to show caller
* **logger:** add report logger in dev mode
* **logger:** add logger struct, modify file log format
* **logger:** add caller func
* **logger:** add LOG_FORMAT to env file
* **naming:** make clean arch naming make sense
* **open:** add limit if open is not equal to last close, then unsubscribe
* **order:** add order cache, add trade agent in stream to decide action, bus as global
* **order:** move order repo from stream usecase to order usecase
* **order:** add group id to recognize parent, remove rsi high low to trade out
* **order:** add realtime data to generate order, order callback to save order
* **order:** add stock update date, add force trade out in trade agent, qty by biasrate
* **order:** add column uuid
* **order:** add all order api
* **pickstock:** add check if stock exist in stream pickstock, add query target order by rank
* **postgres:** add transaction to all repo, add quota, add all orders topic to calculate balance
* **protobuf:** use new format of protobuf, use subscribe future tick to get gap of night market
* **qty:** add modify order qty by return order, trade out order qty will depend on it
* **rabbitmq:** change from grpc stream to rabbitmq, redesign event flow
* **repo:** add postgres index, relation, modify fetching history log, skip close is 0 in fetch
* **repo:** add postgres repo relation, add trade day method
* **rsi:** modify rsi method, add tick time in order, add rsi = 50 as a switch to trade out
* **rsi:** change method calculate rsi, rsi mininum count use as effective count
* **simulate:** temp remove one stock trade once limit and quota, fix simulate difference
* **simulate:** modifiy condition log, default config
* **simulate:** add simulate api, add one trade per stock a day
* **simulate:** replace tickarr to pointer to reduce simulate memory cost
* **sinopac:** implement all sinopac gRPC method
* **status:** add check if in trade time for update order status
* **status:** add fetch history done event to control update order status or not
* **stock:** every time launch, update stock day trade to no, update by latest data
* **stream:** add stream usecase first part, move pb package to outside
* **subscribe:** add search targets on trade day, and subscribe ticks, bidask
* **subscribe:** add subscribe future bidask, modify trade day module
* **target:** add subscribe or not in target cond
* **target:** add black stock, black catagory
* **target:** add alternative choice to find target when volume rank is empty
* **target:** add target cache, add target api, add pgx transaction
* **target:** remove realtime target tag, subscribe first, modify trade in method
* **target:** add realtime rank to config, add debug log in development
* **target:** timer of add realtime target start from 9:00
* **target:** add multiple target condition
* **target:** add target filter struct, modify target cond config
* **target:** remove pre-fetch in target cond, fix wrong target send to analyze
* **target:** add realtime targets, add clear all unfinished orders method
* **target:** add target filter
* **target:** only add target in realtime, last trade day target only use for fetch data
* **ticks:** add fetch history ticks, add grpc max message size to 3G
* **time:** modify trade time unit, add aborted when quota is not enough, waiting will be nil now
* **trade:** modify future agent struct, temp modify future order generator
* **trade:** add future trade agent, temp use the same logic of stock agent
* **trade:** remove trade once limit
* **trade:** use low high compare to all tick out-in-ratio to trade in
* **trade:** change method to analyze volume pr, add use default config to simulate header
* **tradeagent:** add compare with all and period's out in ratio
* **tradein:** add 0.1 all outinratio or inoutratio
* **trader:** move max hold time from global to each agent, add high frequency trade of rsi
* **trader:** refactor future trader, and move to modules, new events pkg
* **trader:** temp remove decide allow forward or reverse by future gap
* **tradeswitch:** add check trade switch every 5 seconds
* **unsubscribe:** add if order canceled, then unsubscribe tick and bidask
* **unsubscribe:** add event topic to unsubscribe(not implement)
* **usecase:** add first usecase, include api, grpc, postgres repo
* **websocket:** add websocket of pickstock on stream router
* **websocket:** add realtime future tick websocket
