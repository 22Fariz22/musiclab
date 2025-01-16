// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Fariz Rustamov",
            "url": "https://github.com/22fariz22",
            "email": "fariz08@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/library": {
            "get": {
                "description": "Возвращает список песен на основе фильтров",
                "tags": [
                    "Songs"
                ],
                "summary": "Получение библиотеки",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Фильтр по группе",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по названию песни",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по тексту",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по дате выпуска",
                        "name": "release_date",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Количество записей на странице",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список песен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/lyrics/ping": {
            "get": {
                "description": "Проверяет доступность базы данных, возвращает \"pong\"",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Проверка доступности базы данных",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/songs": {
            "put": {
                "description": "Обновляет данные песни по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Обновление песни",
                "parameters": [
                    {
                        "description": "Данные для обновления",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateTrackRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Песня успешно обновлена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Создает новую песню на основе данных запроса",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Songs"
                ],
                "summary": "Создание песни",
                "parameters": [
                    {
                        "description": "Данные новой песни",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Созданная песня",
                        "schema": {
                            "$ref": "#/definitions/models.SongDetail"
                        }
                    },
                    "400": {
                        "description": "Некорректные данные"
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера"
                    }
                }
            }
        },
        "/songs/{id}": {
            "delete": {
                "description": "Удаляет песню из базы данных по ID",
                "tags": [
                    "Songs"
                ],
                "summary": "Удаление песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Песня успешно удалена"
                    },
                    "400": {
                        "description": "Некорректный ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/songs/{id}/verses": {
            "get": {
                "description": "Возвращает куплет песни по ID песни и номеру страницы",
                "tags": [
                    "Songs"
                ],
                "summary": "Получение куплета",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Куплет песни",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректный ID или номер страницы",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Куплет не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.SongDetail": {
            "description": "Response containing details of a song",
            "type": "object",
            "required": [
                "link",
                "releaseDate",
                "text"
            ],
            "properties": {
                "link": {
                    "description": "External link to the song\nRequired: true",
                    "type": "string"
                },
                "releaseDate": {
                    "description": "Release date of the song\nRequired: true",
                    "type": "string"
                },
                "text": {
                    "description": "Lyrics or text of the song\nRequired: true",
                    "type": "string"
                }
            }
        },
        "models.SongRequest": {
            "description": "Request payload for adding a new song",
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "description": "Group name of the song\nRequired: true\nMin length: 1",
                    "type": "string",
                    "minLength": 1
                },
                "song": {
                    "description": "Song name\nRequired: true\nMin length: 1",
                    "type": "string",
                    "minLength": 1
                }
            }
        },
        "models.UpdateTrackRequest": {
            "description": "Request payload for updating song details",
            "type": "object",
            "required": [
                "group",
                "id",
                "release_date",
                "song"
            ],
            "properties": {
                "group": {
                    "description": "Group name\nRequired: true\nMin length: 1",
                    "type": "string",
                    "minLength": 1
                },
                "id": {
                    "description": "ID of the track to update\nRequired: true",
                    "type": "integer"
                },
                "link": {
                    "description": "External link to the song",
                    "type": "string"
                },
                "release_date": {
                    "description": "Release date\nRequired: true",
                    "type": "string",
                    "minLength": 1
                },
                "song": {
                    "description": "Song name\nRequired: true\nMin length: 1",
                    "type": "string",
                    "minLength": 1
                },
                "text": {
                    "description": "Lyrics or text of the song",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "MusicLab API",
	Description:      "Это API для управления музыкальной библиотекой.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
