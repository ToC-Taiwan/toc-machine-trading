# [](https://gitlab.tocraw.com/root/toc-machine-trading/compare/v0.0.5...v) (2022-08-08)

## CHANGELOG

## [0.0.5](https://gitlab.tocraw.com/root/toc-machine-trading/compare/v0.0.4...v0.0.5) (2022-08-06)

### Bug Fixes

* **status:** fix fetch history done event wrong input function ([01f980a](https://gitlab.tocraw.com/root/toc-machine-trading/commit/01f980aac99975489d794ee268334e94a7881e99))

### Features

* **event:** remove fetch_history_done event ([b459003](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b459003cbc4fe7dc896b474dd7e698b54ab0ffe7))
* **history:** add fetch history event lock ([4545426](https://gitlab.tocraw.com/root/toc-machine-trading/commit/4545426be2d948e98e3859b17f27dad66fb7541e))
* **status:** add fetch history done event to control update order status or not ([649d848](https://gitlab.tocraw.com/root/toc-machine-trading/commit/649d848b6e55ac8fc1ea1f1ac39808be5a273c0a))
* **target:** remove pre-fetch in target cond, fix wrong target send to analyze ([5ccac3e](https://gitlab.tocraw.com/root/toc-machine-trading/commit/5ccac3ec93314683554454e3e77980c3d7f8301b))

## [0.0.4](https://gitlab.tocraw.com/root/toc-machine-trading/compare/v0.0.3...v0.0.4) (2022-07-29)

### Features

* **config:** remove redundant config, change default config trade fee discount ratio ([600ac2b](https://gitlab.tocraw.com/root/toc-machine-trading/commit/600ac2bf2174b06932d927237a9d00bcda175bb7))
* **stock:** every time launch, update stock day trade to no, update by latest data ([dd1c076](https://gitlab.tocraw.com/root/toc-machine-trading/commit/dd1c076eed2d526908ad2ad3c6160542102cc3ae))
* **target:** add black stock, black catagory ([26c6e41](https://gitlab.tocraw.com/root/toc-machine-trading/commit/26c6e412d972066f6b7afe9867391a1db9b26ff4))

## [0.0.3](https://gitlab.tocraw.com/root/toc-machine-trading/compare/v0.0.2...v0.0.3) (2022-07-28)

### Bug Fixes

* **open:** fix wrong open change ratio in trade ([4d98dad](https://gitlab.tocraw.com/root/toc-machine-trading/commit/4d98dad3b555d7ec0d001f05abce8aad14f5baf8))
* **order:** update order in every order status return, no compare ([ecc786b](https://gitlab.tocraw.com/root/toc-machine-trading/commit/ecc786b58819697f93f16eaefc087032a415727c))
* **target:** fix stuck by non async event, add fetch list in history usecase ([12907d0](https://gitlab.tocraw.com/root/toc-machine-trading/commit/12907d0286828741b593e416df4925cfc3fd4677))
* **target:** fix wrong repo table name when update target ([06798f2](https://gitlab.tocraw.com/root/toc-machine-trading/commit/06798f270a71e6706bd12d4a07a30dc5200c268c))

### Features

* **change-ratio:** consider open change ratio to unsubscribe, and change ratio in stream ([46d9cbd](https://gitlab.tocraw.com/root/toc-machine-trading/commit/46d9cbdfe9fec702055de47da14d34d526ea71a2))
* **qty:** add modify order qty by return order, trade out order qty will depend on it ([a158be9](https://gitlab.tocraw.com/root/toc-machine-trading/commit/a158be9e475588de389655dd307cc806a361c651))
* **target:** add realtime rank to config, add debug log in development ([78eec9a](https://gitlab.tocraw.com/root/toc-machine-trading/commit/78eec9a77b57bb9a4a0cca2961deb34397998995))
* **target:** only add target in realtime, last trade day target only use for fetch data ([d884a21](https://gitlab.tocraw.com/root/toc-machine-trading/commit/d884a214530c1d73dd8a81db1d0e601d2b06814f))
* **target:** timer of add realtime target start from 9:00 ([385d2a9](https://gitlab.tocraw.com/root/toc-machine-trading/commit/385d2a949b1162a2bc89a4b172db39e79de3348b))
* **trade:** change method to analyze volume pr, add use default config to simulate header ([d0a5883](https://gitlab.tocraw.com/root/toc-machine-trading/commit/d0a5883008455ad3b36be3ffc2c65660ccbb820f))
* **unsubscribe:** add if order canceled, then unsubscribe tick and bidask ([b09d65a](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b09d65a88b273fb39cb98d07b81c1f02ffa5de21))

## [0.0.2](https://gitlab.tocraw.com/root/toc-machine-trading/compare/v0.0.1...v0.0.2) (2022-07-21)

### Bug Fixes

* **cancel:** fix cancel fail casue filled order not append to order map, add cancel wait time ([f32a553](https://gitlab.tocraw.com/root/toc-machine-trading/commit/f32a5534566fc436c3e513dec4fc0a37307b16fe))
* **config:** fix simulate must be true ([1aaccb6](https://gitlab.tocraw.com/root/toc-machine-trading/commit/1aaccb6236cfc86523bb51746be91fb4d85d4b28))
* **cpu:** fix checkFirstTickArrive cause cpu 100% ([0f43f43](https://gitlab.tocraw.com/root/toc-machine-trading/commit/0f43f435e663e154947c548ca5d6b442ed54b956))
* **history:** fix skip close is 0 in insert db panic ([b31a256](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b31a256682fdab589c1e1fa2008f993cd046bf99))
* **logger:** fix wrong use format in log error ([70bc20d](https://gitlab.tocraw.com/root/toc-machine-trading/commit/70bc20dfaf90ba12274080e87f1150539e7b3b13))
* **quota:** fix wrong quota when sell or buylater, fix check cancel order wrong tarde time ([40d12fa](https://gitlab.tocraw.com/root/toc-machine-trading/commit/40d12fa27513dd6f9530409e4c5d6a6726adb5e7))
* **snapshot:** fix all snapshot return empty panic, insert all stock from sinopac ([9da88f4](https://gitlab.tocraw.com/root/toc-machine-trading/commit/9da88f4008a848231b0bb2609b2ee55ca43586ec))
* **tradeagent:** fix wrong in out ratio compare ([a48a2e7](https://gitlab.tocraw.com/root/toc-machine-trading/commit/a48a2e755fae16bcc2971a573a367fb91c8aaeac))

### Features

* **analyze:** add last trade day all ticks cahce, upgrade go to 1.18.4 ([d1bc929](https://gitlab.tocraw.com/root/toc-machine-trading/commit/d1bc929a63494c4d2bc29950d8c3b5f332b3ba36))
* **analyze:** same way to analyze realtime tick,history tick, order trade time assign in place ([b5c8a4e](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b5c8a4e028047abd558cc3bb666ba5251c6b9441))
* **basic:** skip add category is 00 to stock list ([b3c847d](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b3c847df328cc1169edd30f7fb85ffd05ed03df6))
* **repo:** add postgres index, relation, modify fetching history log, skip close is 0 in fetch ([e03e2cd](https://gitlab.tocraw.com/root/toc-machine-trading/commit/e03e2cd6733819666932924fe79a8830eb4419b7))
* **simulate:** add simulate api, add one trade per stock a day ([3653954](https://gitlab.tocraw.com/root/toc-machine-trading/commit/3653954c8199bf8e323c367eb270916bd7cddd90))
* **simulate:** replace tickarr to pointer to reduce simulate memory cost ([5964c11](https://gitlab.tocraw.com/root/toc-machine-trading/commit/5964c1114b32b645d8ca959d468fb6987d41e291))
* **status:** add check if in trade time for update order status ([4b92d61](https://gitlab.tocraw.com/root/toc-machine-trading/commit/4b92d61b503215e99f9c05377318316a8c290672))

## [0.0.1](https://gitlab.tocraw.com/root/toc-machine-trading/compare/94fd4a56259d9ca8e7c99417195165358148b2f8...v0.0.1) (2022-07-07)

### Bug Fixes

* **ci:** fix wrong way to load env file ([27de42a](https://gitlab.tocraw.com/root/toc-machine-trading/commit/27de42a131c33b7530b5e072459dd567b6c37c18))
* **event:** fix missing subscribe bidask ([7f28868](https://gitlab.tocraw.com/root/toc-machine-trading/commit/7f28868cf57f513d3868a27f191d59eee9eb96d1))
* **history:** fix wrong key with new fetch kbar tick ([1d28887](https://gitlab.tocraw.com/root/toc-machine-trading/commit/1d2888765c8bbcd93d2f6b655e6970a9e7895306))
* **history:** no gorutine when process data to avoid unexpected error ([c033218](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c03321827761adcc3b2905433ed64732e0a3ecb7))
* **order:** fix empty order time when update, temp extend trade in, out wait time ([9668de5](https://gitlab.tocraw.com/root/toc-machine-trading/commit/9668de54e7563cb9cf0e496218a35928b810dc44))
* **order:** fix repeat place order and cancel ([c785e22](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c785e22993c17146951e0f04bfbc526ec63f84bb))
* **order:** fix wrong out in ratio in order generator ([06ddff7](https://gitlab.tocraw.com/root/toc-machine-trading/commit/06ddff771a91230d875f114252e185707c7efa95))
* **order:** fix wrong status when order updated, add last tick in trader, add lock in cancel order ([b1f307b](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b1f307b366f773ac1b32835dbe0fab827ff60054))
* **order:** using waiting order in tick for, remove uuid, place order fail add status ([949a5b6](https://gitlab.tocraw.com/root/toc-machine-trading/commit/949a5b67f2fc4fc635506d221af9b603800a6164))
* **path:** modify initial sequence, fix wrong event subscribe callback ([0c4499d](https://gitlab.tocraw.com/root/toc-machine-trading/commit/0c4499d831620cd68140a54f068f9250ef56b704))
* **quota:** fix wrong calulate quota, check order status is failed, add lock in place order ([efc4dc5](https://gitlab.tocraw.com/root/toc-machine-trading/commit/efc4dc5bb57bf6d380cace1553eb554d3c3d0c9d))
* **readme:** fix wrong attachment path ([25b2a86](https://gitlab.tocraw.com/root/toc-machine-trading/commit/25b2a86e3c2b19698c6076a6582a48ea8ff12b18))
* **table:** fix wrong table name fail to create in postgres ([8f72bce](https://gitlab.tocraw.com/root/toc-machine-trading/commit/8f72bce677463a570f27e4692a938bb6fec55cd3))
* **trade:** fix wrong analyze tick time, if first tick not arrive no action ([ae4c7c8](https://gitlab.tocraw.com/root/toc-machine-trading/commit/ae4c7c87e19f71c6af29f9e95b6c41a512861481))

### Features

* **agent:** add new agent method, reanme trade to trade agent ([2d809cf](https://gitlab.tocraw.com/root/toc-machine-trading/commit/2d809cffa790894599af28812e5eb2b4ba1dd39a))
* **analyze:** add below quater ma stocks, include api, add history analyze table ([4541f76](https://gitlab.tocraw.com/root/toc-machine-trading/commit/4541f760adf59a6b85d70c202cb8b157327801f4))
* **analyze:** add biasrate of history close ([2aa4c71](https://gitlab.tocraw.com/root/toc-machine-trading/commit/2aa4c71560e9c8c7eb5d479b02a99543d1c2eb85))
* **analyze:** add history open from kbar, fix wrong routing key of bidask, redesign cache ([6e96deb](https://gitlab.tocraw.com/root/toc-machine-trading/commit/6e96deb46eebcd452e893ee0851df648a18a1f82))
* **analyze:** add quater ma to cache ([d0ff1fa](https://gitlab.tocraw.com/root/toc-machine-trading/commit/d0ff1fa8d21934a56a57b78ef6d58173c1c3f956))
* **api:** add day trade calculator, add config api ([c6136d8](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c6136d873417409591690f4a3729d4db0d4514b8))
* **balance:** add all trade balance api ([6f59052](https://gitlab.tocraw.com/root/toc-machine-trading/commit/6f59052a90b421e91dd8ad3cbb1a27e0f95761cf))
* **basic:** add import holiday from data/holiday.json ([c9ec46e](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c9ec46eae8d25d238c004740db2113621a460a92))
* **basic:** add insert or update stock, remove stock entity id, to number ([35a2ca0](https://gitlab.tocraw.com/root/toc-machine-trading/commit/35a2ca09a62d98f125418f4d56a8603c42cf8711))
* **cache:** add stock detail cache in basic usecase ([5ef46a0](https://gitlab.tocraw.com/root/toc-machine-trading/commit/5ef46a0b1111c37123b210356b4b083066ef6a6b))
* **cache:** refactor cache, add basic info, fix insert sql too many parameters ([212d020](https://gitlab.tocraw.com/root/toc-machine-trading/commit/212d02059e4a2b929a9607f0b9b2bb54e06f869b))
* **config:** add sinopac path to config and env, modify readme env template ([42770c5](https://gitlab.tocraw.com/root/toc-machine-trading/commit/42770c5b2b0a22294eb8ce843223984e0f92f950))
* **config:** add tag for config api, rename most api method ([c705195](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c705195518bc81cc7147c0de24112f6c27335c8a))
* **config:** add trade about config, change terminate to put ([6793642](https://gitlab.tocraw.com/root/toc-machine-trading/commit/67936420e3513e194f2914d6c672caa74dd87fc8))
* **event:** add eventbus package, rename from stock to basic ([a66f229](https://gitlab.tocraw.com/root/toc-machine-trading/commit/a66f229ce8fd5a137a613bd34c6485ad80cc7068))
* **event:** add log if event is not about subscribe ([eb958ea](https://gitlab.tocraw.com/root/toc-machine-trading/commit/eb958ea3ca643bd6c5858f972762e3a00fa646d1))
* **grpc:** move snapshot from basic to stream, add tse snapshot api, add snapshot entity ([95ac2bc](https://gitlab.tocraw.com/root/toc-machine-trading/commit/95ac2bc013d37d4e5bc66d7faed3a9b058532ab4))
* **heartbeat:** add heartbeat, history entity modify logger initial ([ed67784](https://gitlab.tocraw.com/root/toc-machine-trading/commit/ed6778478e670383c23d29dc31b93f0cfa483e7c))
* **history:** add day kbar api, add all usecase router ([736f467](https://gitlab.tocraw.com/root/toc-machine-trading/commit/736f467d27ec74631bc4ba202b63a5a24dd2171d))
* **history:** add delete before fetch ([4238a40](https://gitlab.tocraw.com/root/toc-machine-trading/commit/4238a40e362bf6878f6c956809f54a30b2dc0f5f))
* **history:** add history close fetch, add terminate sionopac api, remove open in history close ([5d88611](https://gitlab.tocraw.com/root/toc-machine-trading/commit/5d8861134a1a88075d4f1686bd141dc412854830))
* **history:** add history tick analyze, and use for realtime order generator ([153dd98](https://gitlab.tocraw.com/root/toc-machine-trading/commit/153dd9838e142b0b52f7699c718e257ed4667e3d))
* **history:** modify method to check already exist kbar, tick, merge stream tick, bidask processor ([da32966](https://gitlab.tocraw.com/root/toc-machine-trading/commit/da32966751c4b7cdef2e7680959273a5d2180ae8))
* **kbar:** add fetch history kbar ([17f5ecf](https://gitlab.tocraw.com/root/toc-machine-trading/commit/17f5ecf0ef56dce90adcf393cda961490d81655f))
* **logger:** add caller func ([94fd4a5](https://gitlab.tocraw.com/root/toc-machine-trading/commit/94fd4a56259d9ca8e7c99417195165358148b2f8))
* **naming:** make clean arch naming make sense ([538eb30](https://gitlab.tocraw.com/root/toc-machine-trading/commit/538eb3051f2719f90b288ef004644591be4b0406))
* **order:** add all order api ([c1f6e23](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c1f6e235a0c5b3fe47cba2573e356a12553d096a))
* **order:** add column uuid ([b597adb](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b597adbd92a3fed6c9e31d2cc71a3d7fd6d11107))
* **order:** add order cache, add trade agent in stream to decide action, bus as global ([418ca05](https://gitlab.tocraw.com/root/toc-machine-trading/commit/418ca05e389b34a47543b5a38674d67888c8099d))
* **order:** add realtime data to generate order, order callback to save order ([33ac2b9](https://gitlab.tocraw.com/root/toc-machine-trading/commit/33ac2b9d572a162e2081e3b3d0b3f3e18d08a0c2))
* **order:** add stock update date, add force trade out in trade agent, qty by biasrate ([b14cb16](https://gitlab.tocraw.com/root/toc-machine-trading/commit/b14cb16bcd0c2f2bc73340cc0e27653b7cb31bdd))
* **order:** move order repo from stream usecase to order usecase ([d0f50ec](https://gitlab.tocraw.com/root/toc-machine-trading/commit/d0f50ec7dd20ca80e678a5d7774596d1073a64d7))
* **pickstock:** add check if stock exist in stream pickstock, add query target order by rank ([c132fb4](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c132fb4d99d9273ae83ba185536e7055d3c35b67))
* **postgres:** add transaction to all repo, add quota, add all orders topic to calculate balance ([ef74c35](https://gitlab.tocraw.com/root/toc-machine-trading/commit/ef74c35d41a5487ee4915b53efd890ed27e17fa1))
* **rabbitmq:** change from grpc stream to rabbitmq, redesign event flow ([41033e6](https://gitlab.tocraw.com/root/toc-machine-trading/commit/41033e6a05f69acee6bd612de0c1ee6e21035545))
* **repo:** add postgres repo relation, add trade day method ([16494c8](https://gitlab.tocraw.com/root/toc-machine-trading/commit/16494c82984f3bec9b6ef280b9bd133c8926e56e))
* **sinopac:** implement all sinopac gRPC method ([2b868a9](https://gitlab.tocraw.com/root/toc-machine-trading/commit/2b868a96da2aee5f22c0fd4a2c8f974fcc071159))
* **stream:** add stream usecase first part, move pb package to outside ([240b8ee](https://gitlab.tocraw.com/root/toc-machine-trading/commit/240b8ee886cb76486e78735d2cbf4f815fba9f8f))
* **subscribe:** add search targets on trade day, and subscribe ticks, bidask ([911ef35](https://gitlab.tocraw.com/root/toc-machine-trading/commit/911ef3581007996ec955d745c8ca667d8e4529be))
* **target:** add multiple target condition ([fb8e729](https://gitlab.tocraw.com/root/toc-machine-trading/commit/fb8e7290108e69ae2d00e911ad715ff31bcdeaa9))
* **target:** add realtime targets, add clear all unfinished orders method ([0f8c661](https://gitlab.tocraw.com/root/toc-machine-trading/commit/0f8c6617bbb0c4f0c306d0b5b8e1fb57086f015f))
* **target:** add subscribe or not in target cond ([eb76287](https://gitlab.tocraw.com/root/toc-machine-trading/commit/eb762870aa8039883c4462346d959b730e52fada))
* **target:** add target cache, add target api, add pgx transaction ([94ac9e2](https://gitlab.tocraw.com/root/toc-machine-trading/commit/94ac9e2d897481579ff29e0606f3f346601170cd))
* **target:** add target filter ([49981a4](https://gitlab.tocraw.com/root/toc-machine-trading/commit/49981a4c8ba338c466e17e125df9550a5111e343))
* **ticks:** add fetch history ticks, add grpc max message size to 3G ([36c908e](https://gitlab.tocraw.com/root/toc-machine-trading/commit/36c908eade6f072d616a63123ada4e31c16f5455))
* **time:** modify trade time unit, add aborted when quota is not enough, waiting will be nil now ([de9941c](https://gitlab.tocraw.com/root/toc-machine-trading/commit/de9941cccd412f9dc3d00dac92a40cd1be96014e))
* **tradeswitch:** add check trade switch every 5 seconds ([2e6fda1](https://gitlab.tocraw.com/root/toc-machine-trading/commit/2e6fda1d2e07132c4205a74742c5a6bb9993a731))
* **unsubscribe:** add event topic to unsubscribe(not implement) ([f919d01](https://gitlab.tocraw.com/root/toc-machine-trading/commit/f919d012439a11eacac9a13538557dd278395764))
* **usecase:** add first usecase, include api, grpc, postgres repo ([9800f86](https://gitlab.tocraw.com/root/toc-machine-trading/commit/9800f865e3ed0b16afe06f487523c0a3a4d2b9d8))
* **websocket:** add websocket of pickstock on stream router ([624c83c](https://gitlab.tocraw.com/root/toc-machine-trading/commit/624c83cd9be9ec504b36c9c29ddb12df3aec61e8))
