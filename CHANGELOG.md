<a name="unreleased"></a>
## [Unreleased]

### Chore
- **actions:** modify build and push action version
- **balance:** revert to normal method query future order
- **naming:** reanme needassist to assisting
- **pkg:** move config to module
- **simulation:** add log for user to know is simulation or not
- **stream:** remove out chg and in chg
- **websocket:** revert to new channel in websocket router

### Ci
- **docker:** add git hash to build tag
- **registry:** move registry from docker to github

### Feat
- **assist:** finish by balance and by time period assist
- **assist:** add assist trader in futrue trade ws
- **balance:** add manual trade balance
- **balance:** add manual balance router
- **index:** add otc, tse, nasdaq index to websocket stream
- **order:** remove order status in websocket, modify manual order group id and trade time
- **order:** refactor updateAllTradeBalance, increse update order speed
- **order:** modify new sinopac filling to partfilled
- **order:** remove ask update api, split update balance and simulate order, product order
- **order:** add order status stream in websocket
- **ordertime:** modifiy wrong order time to time now, when in night market from 0:00 to 5:00
- **postition:** add future postion, remove manual in stock, future balance
- **router:** add query future order by date api
- **router:** add manual insert future order api
- **stream:** modify out in volume to four period
- **stream:** add out in rate chg, modify order entity
- **stream:** add period update trade index interal stream usecase
- **stream:** change method process tick arr
- **switch:** add change future trade switch router
- **tick:** cut tick arr in every second in stream
- **trade:** add manual to order column, check order status in websocket
- **trade:** add auto cancel in stream trade
- **websocket:** add nasdaq future to stream

### Fix
- **ci:** add checkout in deployment to get git hash
- **ci:** temp remove go test from ci
- **ci:** add missing config in go test
- **order:** fix order status chan is not add to rabbit, fix order time is always wrong
- **order:** modify get order by trade day will send not filled order
- **order:** use timer and reset to fix balance not insert to db
- **order:** fix order balance calculate wrong, try fix manual order does not insert to db
- **order:** fix manual order does not insert to db
- **order:** add lock for order usecase to update order in postgres
- **position:** fix websocket future position has no column in json
- **postgres:** fix redundant manual future trade balance cause panic
- **router:** fix return body of get future trade switch
- **stream:** add loop lable to avoid index out of range
- **stream:** fix last trade rate not initial
- **stream:** fix map is not initail
- **stuck:** fix missing go in updateAllTradeBalance
- **trade:** fix order will be cancel multiple times
- **trade:** fix order will be cancel before 10 seconds
- **websocket:** fix send data to close channel
- **websocket:** fix order map does not initial
- **ws:** fix ws end abnormal, close connection before gin done

### Refactor
- **logger:** pack log again
- **ws:** split pick stock and future to different pkg
- **ws:** refactor websocket split pick stock and future trade

### Revert
- **stream:** revert to cut and process in the same loop


<a name="v1.5.0"></a>
## [v1.5.0] - 2022-11-16
### Chore
- **entity:** change db table name to split stock future, modify target stream entity
- **entity:** split trade balance to stock and future
- **entity:** modify history entity to split history data and base
- **global:** remove global pkg, move to common pkg
- **log:** remove redundant log of websocket
- **websocket:** add log for unsupport message, remove if v == pong

### Feat
- **action:** use period tick out in ratio to decide action, and add to simulator
- **history:** add future history close fetch and simulate
- **holiday:** extend trade year to 2023, update holiday.json
- **order:** modify all order router to future and stock all order
- **tick:** modify future tick chan and connection id
- **trade:** add ws trade
- **tradeperiod:** add get last 1 trade period for future method
- **websocket:** modify trade rate content to out, in and period
- **websocket:** add trade rate in websocket
- **websocket:** ignore CloseNoStatusReceived error
- **ws:** add log for new future ws and done log
- **ws:** modify ws layout, add send snapshot in future stream ws

### Fix
- **health:** if disconnect from grpc not panic but os exit
- **migration:** fix wrong table name in migration sql
- **subscribe:** remove redundant future bidask subscribe
- **websocket:** remove period to fix out of index
- **websocket:** add missing socketPickStock
- **websocket:** fix interface cast bug
- **websocket:** fix concurrency write websocket
- **websocket:** fix wrong scoket data type
- **websocket:** fix wrong calculation of trade rate
- **websocket:** fix missing format in send websocket data


<a name="v1.4.0"></a>
## [v1.4.0] - 2022-11-06
### Feat
- **future:** use rsi to decide trade in, max hold time to trade out
- **future:** test new method of future trader
- **future:** modify trade out method
- **future:** add kbar analyze to trade out of future
- **router:** add simulate future to history router
- **simulate:** add simulate future trade proto type
- **strategy:** modify future trade strategy default config
- **trader:** add kbar analyze to future trader


<a name="v1.3.0"></a>
## [v1.3.0] - 2022-10-24
### Chore
- **changelog:** modify changelog
- **changelog:** add git-chglog instead of cz-conventional-changelog
- **make:** copy default config in every run, modify stock trader

### Feat
- **balance:** exclude unfinish order in calculate balance, modify topic name
- **router:** add future trade balance to get balance
- **strategy:** modify future trader strategy


<a name="v1.2.0"></a>
## [v1.2.0] - 2022-10-23
### Chore
- **readme:** move tools to contributing doc, add install migrate in go update

### Refactor
- **trader:** put stock future trader into a module, change the stream usecase


<a name="v1.1.0"></a>
## [v1.1.0] - 2022-10-23
### Chore
- **module:** move cache to a new module
- **trader:** remove redundant bus, rename trade agent to trader

### Feat
- **future:** add auto select r1 mxf future
- **module:** move event topic to a new event module with bus
- **module:** move target filter to a new target module
- **module:** move all trader to trader module, remove simulate stock trader
- **tradeday:** add new future trade day method, modify query order method to time range

### Fix
- **future:** fix wrong future trade switch, add future trade fee calculate


<a name="v1.0.0"></a>
## v1.0.0 - 2022-10-23
### Chore
- **bidask:** change tick time to bidask time
- **cache:** rename cache method
- **config:** modifiy default config
- **dependency:** update go dependency, protoc version
- **entity:** move order status map to order entity
- **event:** change event name from target to fetch history
- **format:** format all import, add global time format without dash
- **future:** add new analyze method in future trade
- **go:** update go version to 1.19
- **log:** modify usage of logger
- **logger:** remove reporter in logrus
- **mod:** change go mod name to tmt
- **module:** move trade to a new module
- **module:** add quota module
- **naming:** rename stock simulate trader file name, add future simulate file
- **naming:** rename realtime data to trader
- **port:** change default http port
- **protobuf:** modify protobuf folder level
- **readme:** remove old url, add badge in reademe
- **reorder:** reorder analyze usecase
- **stream:** remove period get tse snapshot in stream usecase

### Ci
- **ci:** add dockerfile and gitlab-ci
- **env:** reset env and config after kill container
- **migrate:** migrate from gitlab ci to github actions
- **port:** add machine port in actions

### Docs
- **changelog:** v0.0.4
- **changelog:** modify changelog
- **readme:** add clean arch layers image
- **script:** modify makefile style, add pre commit, remove callvis install

### Feat
- **agent:** add new agent method, reanme trade to trade agent
- **analyze:** split stock, future analyze, fix future order repo, future first tick analyze time
- **analyze:** add last trade day all ticks cahce, upgrade go to 1.18.4
- **analyze:** add quater ma to cache
- **analyze:** add biasrate of history close
- **analyze:** add history open from kbar, fix wrong routing key of bidask, redesign cache
- **analyze:** add below quater ma stocks, include api, add history analyze table
- **analyze:** same way to analyze realtime tick,history tick, order trade time assign in place
- **api:** add day trade calculator, add config api
- **balance:** add all trade balance api
- **balance:** add calculate future order balance, add allow trade in stock, future config
- **basic:** skip add category is 00 to stock list
- **basic:** add import holiday from data/holiday.json
- **basic:** add insert or update stock, remove stock entity id, to number
- **biasrate:** change the way of biasRate usage, use bidask price to decide trade out price
- **cache:** add stock detail cache in basic usecase
- **cache:** refactor cache, add basic info, fix insert sql too many parameters
- **change-ratio:** consider open change ratio to unsubscribe, and change ratio in stream
- **config:** add sinopac path to config and env, modify readme env template
- **config:** add trade about config, change terminate to put
- **config:** remove redundant config, change default config trade fee discount ratio
- **config:** add tag for config api, rename most api method
- **entity:** add base order to distinguish stock order and future order
- **event:** add eventbus package, rename from stock to basic
- **event:** add log if event is not about subscribe
- **event:** remove fetch_history_done event
- **future:** add future detail, add subscribe future tick
- **future:** use future gap in 8:45 to decide forward or reverse, add trade day struct
- **go:** update go to 1.19.1 and remove k8s agent
- **grpc:** move snapshot from basic to stream, add tse snapshot api, add snapshot entity
- **heartbeat:** add heartbeat, history entity modify logger initial
- **history:** add fetch history event lock
- **history:** add history tick analyze, and use for realtime order generator
- **history:** add delete before fetch
- **history:** modify method to check already exist kbar, tick, merge stream tick, bidask processor
- **history:** add history close fetch, add terminate sionopac api, remove open in history close
- **history:** add day kbar api, add all usecase router
- **kbar:** add fetch history kbar
- **log:** merge all log into one file, and pretty json
- **logger:** add report logger in dev mode
- **logger:** add LOG_FORMAT to env file
- **logger:** add logger struct, modify file log format
- **logger:** modify log format and check if development to show caller
- **logger:** add caller func
- **naming:** make clean arch naming make sense
- **open:** add limit if open is not equal to last close, then unsubscribe
- **order:** add column uuid
- **order:** add order cache, add trade agent in stream to decide action, bus as global
- **order:** move order repo from stream usecase to order usecase
- **order:** add all order api
- **order:** add group id to recognize parent, remove rsi high low to trade out
- **order:** add realtime data to generate order, order callback to save order
- **order:** add stock update date, add force trade out in trade agent, qty by biasrate
- **pickstock:** add check if stock exist in stream pickstock, add query target order by rank
- **postgres:** add transaction to all repo, add quota, add all orders topic to calculate balance
- **protobuf:** use new format of protobuf, use subscribe future tick to get gap of night market
- **qty:** add modify order qty by return order, trade out order qty will depend on it
- **rabbitmq:** change from grpc stream to rabbitmq, redesign event flow
- **repo:** add postgres index, relation, modify fetching history log, skip close is 0 in fetch
- **repo:** add postgres repo relation, add trade day method
- **rsi:** modify rsi method, add tick time in order, add rsi = 50 as a switch to trade out
- **rsi:** change method calculate rsi, rsi mininum count use as effective count
- **simulate:** replace tickarr to pointer to reduce simulate memory cost
- **simulate:** temp remove one stock trade once limit and quota, fix simulate difference
- **simulate:** modifiy condition log, default config
- **simulate:** add simulate api, add one trade per stock a day
- **sinopac:** implement all sinopac gRPC method
- **status:** add fetch history done event to control update order status or not
- **status:** add check if in trade time for update order status
- **stock:** every time launch, update stock day trade to no, update by latest data
- **stream:** add stream usecase first part, move pb package to outside
- **subscribe:** add search targets on trade day, and subscribe ticks, bidask
- **subscribe:** add subscribe future bidask, modify trade day module
- **target:** add subscribe or not in target cond
- **target:** add target cache, add target api, add pgx transaction
- **target:** remove pre-fetch in target cond, fix wrong target send to analyze
- **target:** add realtime rank to config, add debug log in development
- **target:** remove realtime target tag, subscribe first, modify trade in method
- **target:** only add target in realtime, last trade day target only use for fetch data
- **target:** add multiple target condition
- **target:** add realtime targets, add clear all unfinished orders method
- **target:** timer of add realtime target start from 9:00
- **target:** add alternative choice to find target when volume rank is empty
- **target:** add target filter
- **target:** add target filter struct, modify target cond config
- **target:** add black stock, black catagory
- **ticks:** add fetch history ticks, add grpc max message size to 3G
- **time:** modify trade time unit, add aborted when quota is not enough, waiting will be nil now
- **trade:** change method to analyze volume pr, add use default config to simulate header
- **trade:** remove trade once limit
- **trade:** add future trade agent, temp use the same logic of stock agent
- **trade:** modify future agent struct, temp modify future order generator
- **trade:** use low high compare to all tick out-in-ratio to trade in
- **tradeagent:** add compare with all and period's out in ratio
- **tradein:** add 0.1 all outinratio or inoutratio
- **trader:** refactor future trader, and move to modules, new events pkg
- **trader:** move max hold time from global to each agent, add high frequency trade of rsi
- **trader:** temp remove decide allow forward or reverse by future gap
- **tradeswitch:** add check trade switch every 5 seconds
- **unsubscribe:** add if order canceled, then unsubscribe tick and bidask
- **unsubscribe:** add event topic to unsubscribe(not implement)
- **usecase:** add first usecase, include api, grpc, postgres repo
- **websocket:** add realtime future tick websocket
- **websocket:** add websocket of pickstock on stream router

### Fix
- **actions:** add reset environment before deployment
- **cancel:** fix cancel fail casue filled order not append to order map, add cancel wait time
- **ci:** fix wrong way to load env file
- **config:** fix simulate must be true
- **cpu:** fix checkFirstTickArrive cause cpu 100%
- **deployment:** fix wrong config path in deployment action
- **event:** fix missing subscribe bidask
- **future:** fix wrong rsi gap base, rename future trader and stock trader
- **history:** no gorutine when process data to avoid unexpected error
- **history:** fix wrong key with new fetch kbar tick
- **history:** fix skip close is 0 in insert db panic
- **logger:** fix debug missing format function
- **logger:** fix wrong use format in log error
- **open:** fix wrong open change ratio in trade
- **order:** fix repeat place order and cancel
- **order:** fix wrong out in ratio in order generator
- **order:** update order in every order status return, no compare
- **order:** using waiting order in tick for, remove uuid, place order fail add status
- **order:** fix empty order time when update, temp extend trade in, out wait time
- **order:** fix wrong status when order updated, add last tick in trader, add lock in cancel order
- **path:** modify initial sequence, fix wrong event subscribe callback
- **quota:** fix wrong calulate quota, check order status is failed, add lock in place order
- **quota:** fix wrong quota when sell or buylater, fix check cancel order wrong tarde time
- **readme:** fix wrong attachment path
- **simulate:** fix wrong orders return, modify default config, add realtime target to 30 secs
- **snapshot:** fix all snapshot return empty panic, insert all stock from sinopac
- **status:** fix fetch history done event wrong input function
- **table:** fix wrong table name fail to create in postgres
- **target:** fix wrong repo table name when update target
- **target:** fix stuck by non async event, add fetch list in history usecase
- **trade:** fix wrong analyze tick time, if first tick not arrive no action
- **tradeagent:** fix wrong in out ratio compare

### Refactor
- **config:** move config pkg to top level between cmd
- **dependency:** remove logger dependency in all pkg
- **logger:** remove global basepath set and get, remove global dependency from logger
- **pkg:** rename sinopac to grpc, move global to top level
- **pkg:** refactor config, eventbus method
- **target:** rename subscribe and realtime add to pre-fetch, realtime

### Style
- **clean:** first commit from clean code layout


[Unreleased]: https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.5.0...HEAD
[v1.5.0]: https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.4.0...v1.5.0
[v1.4.0]: https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.3.0...v1.4.0
[v1.3.0]: https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.2.0...v1.3.0
[v1.2.0]: https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/ToC-Taiwan/toc-machine-trading/compare/v1.0.0...v1.1.0
