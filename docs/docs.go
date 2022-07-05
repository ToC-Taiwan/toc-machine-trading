// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/analyze/reborn": {
            "get": {
                "description": "getRebornTargets",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "analyze"
                ],
                "summary": "getRebornTargets",
                "operationId": "getRebornTargets",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/v1.reborn"
                            }
                        }
                    }
                }
            }
        },
        "/basic/config": {
            "get": {
                "description": "getAllConfig",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "getAllConfig",
                "operationId": "getAllConfig",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/config.Config"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/basic/stock": {
            "get": {
                "description": "getAllRepoStock",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "basic"
                ],
                "summary": "getAllRepoStock",
                "operationId": "getAllRepoStock",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.stockDetailResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/basic/stock/sinopac-to-repo": {
            "get": {
                "description": "getAllSinopacStockAndUpdateRepo",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "basic"
                ],
                "summary": "getAllSinopacStockAndUpdateRepo",
                "operationId": "getAllSinopacStockAndUpdateRepo",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.stockDetailResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/basic/system/terminate": {
            "put": {
                "description": "terminateSinopac",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "terminateSinopac",
                "operationId": "terminateSinopac",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/history/day-kbar/{stock}/{start_date}/{interval}": {
            "get": {
                "description": "getKbarData",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "history"
                ],
                "summary": "getKbarData",
                "operationId": "getKbarData",
                "parameters": [
                    {
                        "type": "string",
                        "description": "stock",
                        "name": "stock",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "start_date",
                        "name": "start_date",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "interval",
                        "name": "interval",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.HistoryKbar"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/order/all": {
            "get": {
                "description": "getAllOrder",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "getAllOrder",
                "operationId": "getAllOrder",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Order"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/order/balance": {
            "get": {
                "description": "getAllTradeBalance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "getAllTradeBalance",
                "operationId": "getAllTradeBalance",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.TradeBalance"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/order/day-trade/forward": {
            "get": {
                "description": "calculateForwardDayTradeBalance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "calculateForwardDayTradeBalance",
                "operationId": "calculateForwardDayTradeBalance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "buy_price",
                        "name": "buy_price",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "buy_quantity",
                        "name": "buy_quantity",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "sell_price",
                        "name": "sell_price",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "sell_quantity",
                        "name": "sell_quantity",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.dayTradeResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/order/day-trade/reverse": {
            "get": {
                "description": "calculateReverseDayTradeBalance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "calculateReverseDayTradeBalance",
                "operationId": "calculateReverseDayTradeBalance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "sell_first_price",
                        "name": "sell_first_price",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "sell_first_quantity",
                        "name": "sell_first_quantity",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "buy_later_price",
                        "name": "buy_later_price",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "buy_later_quantity",
                        "name": "buy_later_quantity",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.dayTradeResult"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/stream/tse/snapshot": {
            "get": {
                "description": "getTSESnapshot",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stream"
                ],
                "summary": "getTSESnapshot",
                "operationId": "getTSESnapshot",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.StockSnapShot"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        },
        "/targets": {
            "get": {
                "description": "getTargets",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "targets"
                ],
                "summary": "getTargets",
                "operationId": "getTargets",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Target"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "config.Analyze": {
            "type": "object",
            "properties": {
                "close_change_ratio_high": {
                    "type": "number"
                },
                "close_change_ratio_low": {
                    "type": "number"
                },
                "in_out_ratio": {
                    "type": "number"
                },
                "ma_period": {
                    "type": "integer"
                },
                "max_loss": {
                    "type": "number"
                },
                "open_close_change_ratio_high": {
                    "type": "number"
                },
                "open_close_change_ratio_low": {
                    "type": "number"
                },
                "out_in_ratio": {
                    "type": "number"
                },
                "rsi_high": {
                    "type": "number"
                },
                "rsi_low": {
                    "type": "number"
                },
                "rsi_min_count": {
                    "type": "integer"
                },
                "tick_analyze_max_period": {
                    "type": "number"
                },
                "tick_analyze_min_period": {
                    "type": "number"
                },
                "volume_pr_high": {
                    "type": "number"
                },
                "volume_pr_low": {
                    "type": "number"
                }
            }
        },
        "config.Config": {
            "type": "object",
            "properties": {
                "analyze": {
                    "$ref": "#/definitions/config.Analyze"
                },
                "deployment": {
                    "type": "string"
                },
                "history": {
                    "$ref": "#/definitions/config.History"
                },
                "http": {
                    "$ref": "#/definitions/config.HTTP"
                },
                "postgres": {
                    "$ref": "#/definitions/config.Postgres"
                },
                "quota": {
                    "$ref": "#/definitions/config.Quota"
                },
                "rabbitmq": {
                    "$ref": "#/definitions/config.RabbitMQ"
                },
                "sinopac": {
                    "$ref": "#/definitions/config.Sinopac"
                },
                "target_cond": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/config.TargetCond"
                    }
                },
                "trade_switch": {
                    "$ref": "#/definitions/config.TradeSwitch"
                }
            }
        },
        "config.HTTP": {
            "type": "object",
            "properties": {
                "port": {
                    "type": "string"
                }
            }
        },
        "config.History": {
            "type": "object",
            "properties": {
                "history_close_period": {
                    "type": "integer"
                },
                "history_kbar_period": {
                    "type": "integer"
                },
                "history_tick_period": {
                    "type": "integer"
                }
            }
        },
        "config.Postgres": {
            "type": "object",
            "properties": {
                "db_name": {
                    "type": "string"
                },
                "pool_max": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "config.Quota": {
            "type": "object",
            "properties": {
                "fee_discount": {
                    "type": "number"
                },
                "trade_fee_ratio": {
                    "type": "number"
                },
                "trade_quota": {
                    "type": "integer"
                },
                "trade_tax_ratio": {
                    "type": "number"
                }
            }
        },
        "config.RabbitMQ": {
            "type": "object",
            "properties": {
                "attempts": {
                    "type": "integer"
                },
                "exchange": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "wait_time": {
                    "type": "integer"
                }
            }
        },
        "config.Sinopac": {
            "type": "object",
            "properties": {
                "pool_max": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "config.TargetCond": {
            "type": "object",
            "properties": {
                "limit_price_high": {
                    "type": "number"
                },
                "limit_price_low": {
                    "type": "number"
                },
                "limit_volume": {
                    "type": "integer"
                },
                "subscribe": {
                    "type": "boolean"
                }
            }
        },
        "config.TradeSwitch": {
            "type": "object",
            "properties": {
                "buy": {
                    "type": "boolean"
                },
                "buy_later": {
                    "type": "boolean"
                },
                "forward_max": {
                    "type": "integer"
                },
                "hold_time_from_open": {
                    "type": "number"
                },
                "mean_time_forward": {
                    "type": "integer"
                },
                "mean_time_reverse": {
                    "type": "integer"
                },
                "reverse_max": {
                    "type": "integer"
                },
                "sell": {
                    "type": "boolean"
                },
                "sell_first": {
                    "type": "boolean"
                },
                "simulation": {
                    "type": "boolean"
                },
                "total_open_time": {
                    "type": "number"
                },
                "trade_in_end_time": {
                    "type": "number"
                },
                "trade_in_wait_time": {
                    "type": "integer"
                },
                "trade_out_end_time": {
                    "type": "number"
                },
                "trade_out_wait_time": {
                    "type": "integer"
                }
            }
        },
        "entity.HistoryKbar": {
            "type": "object",
            "properties": {
                "close": {
                    "type": "number"
                },
                "high": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "kbar_time": {
                    "type": "string"
                },
                "low": {
                    "type": "number"
                },
                "open": {
                    "type": "number"
                },
                "stock": {
                    "$ref": "#/definitions/entity.Stock"
                },
                "stock_num": {
                    "type": "string"
                },
                "volume": {
                    "type": "integer"
                }
            }
        },
        "entity.Order": {
            "type": "object",
            "properties": {
                "action": {
                    "type": "integer"
                },
                "order_id": {
                    "type": "string"
                },
                "order_time": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "stock": {
                    "$ref": "#/definitions/entity.Stock"
                },
                "stock_num": {
                    "type": "string"
                },
                "trade_time": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "entity.Stock": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "string"
                },
                "day_trade": {
                    "type": "boolean"
                },
                "exchange": {
                    "type": "string"
                },
                "last_close": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "number": {
                    "type": "string"
                },
                "update_date": {
                    "type": "string"
                }
            }
        },
        "entity.StockSnapShot": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "amount_sum": {
                    "type": "integer"
                },
                "chg_type": {
                    "type": "string"
                },
                "close": {
                    "type": "number"
                },
                "high": {
                    "type": "number"
                },
                "low": {
                    "type": "number"
                },
                "open": {
                    "type": "number"
                },
                "pct_chg": {
                    "type": "number"
                },
                "price_chg": {
                    "type": "number"
                },
                "snap_time": {
                    "type": "string"
                },
                "stock_name": {
                    "type": "string"
                },
                "stock_num": {
                    "type": "string"
                },
                "tick_type": {
                    "type": "string"
                },
                "volume": {
                    "type": "integer"
                },
                "volume_ratio": {
                    "type": "number"
                },
                "volume_sum": {
                    "type": "integer"
                },
                "yesterday_volume": {
                    "type": "number"
                }
            }
        },
        "entity.Target": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "rank": {
                    "type": "integer"
                },
                "real_time_add": {
                    "type": "boolean"
                },
                "stock": {
                    "$ref": "#/definitions/entity.Stock"
                },
                "stock_num": {
                    "type": "string"
                },
                "subscribe": {
                    "type": "boolean"
                },
                "trade_day": {
                    "type": "string"
                },
                "volume": {
                    "type": "integer"
                }
            }
        },
        "entity.TradeBalance": {
            "type": "object",
            "properties": {
                "discount": {
                    "type": "integer"
                },
                "forward": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "original_balance": {
                    "type": "integer"
                },
                "reverse": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                },
                "trade_count": {
                    "type": "integer"
                },
                "trade_day": {
                    "type": "string"
                }
            }
        },
        "v1.dayTradeResult": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                }
            }
        },
        "v1.reborn": {
            "type": "object",
            "properties": {
                "date": {
                    "type": "string"
                },
                "stocks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Stock"
                    }
                }
            }
        },
        "v1.response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "v1.stockDetailResponse": {
            "type": "object",
            "properties": {
                "stock_detail": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Stock"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "TOC MACHINE TRADING",
	Description:      "API docs for Auto Trade",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
