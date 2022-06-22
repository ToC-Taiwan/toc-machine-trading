# (2022-06-22)

## CHANGELOG

### Bug Fixes

* **event:** fix missing subscribe bidask ([7f28868](https://gitlab.tocraw.com/root/toc-machine-trading/commit/7f28868cf57f513d3868a27f191d59eee9eb96d1))
* **path:** modify initial sequence, fix wrong event subscribe callback ([0c4499d](https://gitlab.tocraw.com/root/toc-machine-trading/commit/0c4499d831620cd68140a54f068f9250ef56b704))
* **readme:** fix wrong attachment path ([25b2a86](https://gitlab.tocraw.com/root/toc-machine-trading/commit/25b2a86e3c2b19698c6076a6582a48ea8ff12b18))
* **table:** fix wrong table name fail to create in postgres ([8f72bce](https://gitlab.tocraw.com/root/toc-machine-trading/commit/8f72bce677463a570f27e4692a938bb6fec55cd3))

### Features

* **basic:** add import holiday from data/holiday.json ([c9ec46e](https://gitlab.tocraw.com/root/toc-machine-trading/commit/c9ec46eae8d25d238c004740db2113621a460a92))
* **basic:** add insert or update stock, remove stock entity id, to number ([35a2ca0](https://gitlab.tocraw.com/root/toc-machine-trading/commit/35a2ca09a62d98f125418f4d56a8603c42cf8711))
* **cache:** add stock detail cache in basic usecase ([5ef46a0](https://gitlab.tocraw.com/root/toc-machine-trading/commit/5ef46a0b1111c37123b210356b4b083066ef6a6b))
* **cache:** refactor cache, add basic info, fix insert sql too many parameters ([212d020](https://gitlab.tocraw.com/root/toc-machine-trading/commit/212d02059e4a2b929a9607f0b9b2bb54e06f869b))
* **config:** add sinopac path to config and env, modify readme env template ([42770c5](https://gitlab.tocraw.com/root/toc-machine-trading/commit/42770c5b2b0a22294eb8ce843223984e0f92f950))
* **config:** add trade about config, change terminate to put ([6793642](https://gitlab.tocraw.com/root/toc-machine-trading/commit/67936420e3513e194f2914d6c672caa74dd87fc8))
* **event:** add eventbus package, rename from stock to basic ([a66f229](https://gitlab.tocraw.com/root/toc-machine-trading/commit/a66f229ce8fd5a137a613bd34c6485ad80cc7068))
* **heartbeat:** add heartbeat, history entity modify logger initial ([ed67784](https://gitlab.tocraw.com/root/toc-machine-trading/commit/ed6778478e670383c23d29dc31b93f0cfa483e7c))
* **history:** add history close fetch, add terminate sionopac api, remove open in history close ([5d88611](https://gitlab.tocraw.com/root/toc-machine-trading/commit/5d8861134a1a88075d4f1686bd141dc412854830))
* **history:** modify method to check already exist kbar, tick, merge stream tick, bidask processor ([da32966](https://gitlab.tocraw.com/root/toc-machine-trading/commit/da32966751c4b7cdef2e7680959273a5d2180ae8))
* **kbar:** add fetch history kbar ([17f5ecf](https://gitlab.tocraw.com/root/toc-machine-trading/commit/17f5ecf0ef56dce90adcf393cda961490d81655f))
* **logger:** add caller func ([94fd4a5](https://gitlab.tocraw.com/root/toc-machine-trading/commit/94fd4a56259d9ca8e7c99417195165358148b2f8))
* **naming:** make clean arch naming make sense ([538eb30](https://gitlab.tocraw.com/root/toc-machine-trading/commit/538eb3051f2719f90b288ef004644591be4b0406))
* **order:** add order cache, add trade agent in stream to decide action, bus as global ([418ca05](https://gitlab.tocraw.com/root/toc-machine-trading/commit/418ca05e389b34a47543b5a38674d67888c8099d))
* **rabbitmq:** change from grpc stream to rabbitmq, redesign event flow ([41033e6](https://gitlab.tocraw.com/root/toc-machine-trading/commit/41033e6a05f69acee6bd612de0c1ee6e21035545))
* **repo:** add postgres repo relation, add trade day method ([16494c8](https://gitlab.tocraw.com/root/toc-machine-trading/commit/16494c82984f3bec9b6ef280b9bd133c8926e56e))
* **sinopac:** implement all sinopac gRPC method ([2b868a9](https://gitlab.tocraw.com/root/toc-machine-trading/commit/2b868a96da2aee5f22c0fd4a2c8f974fcc071159))
* **stream:** add stream usecase first part, move pb package to outside ([240b8ee](https://gitlab.tocraw.com/root/toc-machine-trading/commit/240b8ee886cb76486e78735d2cbf4f815fba9f8f))
* **subscribe:** add search targets on trade day, and subscribe ticks, bidask ([911ef35](https://gitlab.tocraw.com/root/toc-machine-trading/commit/911ef3581007996ec955d745c8ca667d8e4529be))
* **target:** add target filter ([49981a4](https://gitlab.tocraw.com/root/toc-machine-trading/commit/49981a4c8ba338c466e17e125df9550a5111e343))
* **ticks:** add fetch history ticks, add grpc max message size to 3G ([36c908e](https://gitlab.tocraw.com/root/toc-machine-trading/commit/36c908eade6f072d616a63123ada4e31c16f5455))
* **usecase:** add first usecase, include api, grpc, postgres repo ([9800f86](https://gitlab.tocraw.com/root/toc-machine-trading/commit/9800f865e3ed0b16afe06f487523c0a3a4d2b9d8))
