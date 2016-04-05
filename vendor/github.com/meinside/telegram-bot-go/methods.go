// https://core.telegram.org/bots/api#available-methods
package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Get info of this bot.
//
// https://core.telegram.org/bots/api#getme
func (b *Bot) GetMe() (result ApiResultUser) {
	return b.requestResultUser("getMe", map[string]interface{}{}) // no params
}

// Send a message.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: parse_mode, disable_web_page_preview, reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendmessage
func (b *Bot) SendMessage(chatId interface{}, text *string, options map[string]interface{}) (result ApiResultMessage) {
	// essential params
	params := map[string]interface{}{
		"chat_id": chatId,
		"text":    *text,
	}
	// optional params
	for key, val := range options {
		if val != nil {
			params[key] = val
		}
	}

	return b.requestResultMessage("sendMessage", params)
}

// Forward a message.
//
// chatId and fromChatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// https://core.telegram.org/bots/api#forwardmessage
func (b *Bot) ForwardMessage(chatId interface{}, fromChatId interface{}, messageId int) (result ApiResultMessage) {
	// essential params
	params := map[string]interface{}{
		"chat_id":      chatId,
		"from_chat_id": fromChatId,
		"message_id":   messageId,
	}

	return b.requestResultMessage("forwardMessage", params)
}

// Send photos.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: caption, reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendphoto
func (b *Bot) SendPhoto(chatId interface{}, photoFilepath *string, options map[string]interface{}) (result ApiResultMessage) {
	if file, err := os.Open(*photoFilepath); err == nil {
		// essential params
		params := map[string]interface{}{
			"chat_id": chatId,
			"photo":   file,
		}
		// optional params
		for key, val := range options {
			if val != nil {
				params[key] = val
			}
		}

		return b.requestResultMessage("sendPhoto", params)
	} else {
		errStr := err.Error()

		b.error(errStr)

		return ApiResultMessage{Ok: false, Description: &errStr}
	}
}

// Send audio files. (.mp3 format only, will be played with external players)
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: duration, performer, title, reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendaudio
func (b *Bot) SendAudio(chatId interface{}, audioFilepath *string, options map[string]interface{}) (result ApiResultMessage) {
	if file, err := os.Open(*audioFilepath); err == nil {
		// essential params
		params := map[string]interface{}{
			"chat_id": chatId,
			"audio":   file,
		}
		// optional params
		for key, val := range options {
			if val != nil {
				params[key] = val
			}
		}

		return b.requestResultMessage("sendAudio", params)
	} else {
		errStr := err.Error()

		b.error(errStr)

		return ApiResultMessage{Ok: false, Description: &errStr}
	}
}

// Send general files.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#senddocument
func (b *Bot) SendDocument(chatId interface{}, documentFilepath *string, options map[string]interface{}) (result ApiResultMessage) {
	if file, err := os.Open(*documentFilepath); err == nil {
		// essential params
		params := map[string]interface{}{
			"chat_id":  chatId,
			"document": file,
		}
		// optional params
		for key, val := range options {
			if val != nil {
				params[key] = val
			}
		}

		return b.requestResultMessage("sendDocument", params)
	} else {
		errStr := err.Error()

		b.error(errStr)

		return ApiResultMessage{Ok: false, Description: &errStr}
	}
}

// Send stickers.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendsticker
func (b *Bot) SendSticker(chatId interface{}, stickerFilepath *string, options map[string]interface{}) (result ApiResultMessage) {
	if file, err := os.Open(*stickerFilepath); err == nil {
		// essential params
		params := map[string]interface{}{
			"chat_id": chatId,
			"sticker": file,
		}
		// optional params
		for key, val := range options {
			if val != nil {
				params[key] = val
			}
		}

		return b.requestResultMessage("sendSticker", params)
	} else {
		errStr := err.Error()

		b.error(errStr)

		return ApiResultMessage{Ok: false, Description: &errStr}
	}
}

// Send video files.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: duration, caption, reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendvideo
func (b *Bot) SendVideo(chatId interface{}, videoFilepath *string, options map[string]interface{}) (result ApiResultMessage) {
	if file, err := os.Open(*videoFilepath); err == nil {
		// essential params
		params := map[string]interface{}{
			"chat_id": chatId,
			"video":   file,
		}
		// optional params
		for key, val := range options {
			if val != nil {
				params[key] = val
			}
		}

		return b.requestResultMessage("sendVideo", params)
	} else {
		errStr := err.Error()

		b.error(errStr)

		return ApiResultMessage{Ok: false, Description: &errStr}
	}
}

// Send voice files. (.ogg format only, will be played with Telegram itself))
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: duration, reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendvoice
func (b *Bot) SendVoice(chatId interface{}, voiceFilepath *string, options map[string]interface{}) (result ApiResultMessage) {
	if file, err := os.Open(*voiceFilepath); err == nil {
		// essential params
		params := map[string]interface{}{
			"chat_id": chatId,
			"voice":   file,
		}
		// optional params
		for key, val := range options {
			if val != nil {
				params[key] = val
			}
		}

		return b.requestResultMessage("sendVoice", params)
	} else {
		errStr := err.Error()

		b.error(errStr)

		return ApiResultMessage{Ok: false, Description: &errStr}
	}
}

// Send locations.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// options include: reply_to_message_id, and reply_markup.
//
// https://core.telegram.org/bots/api#sendlocation
func (b *Bot) SendLocation(chatId interface{}, latitude, longitude float32, options map[string]interface{}) (result ApiResultMessage) {
	// essential params
	params := map[string]interface{}{
		"chat_id":   chatId,
		"latitude":  latitude,
		"longitude": longitude,
	}
	// optional params
	for key, val := range options {
		if val != nil {
			params[key] = val
		}
	}

	return b.requestResultMessage("sendLocation", params)
}

// Send chat action.
//
// chatId can be Message.Chat.Id or target channel(eg. @channelusername).
//
// action can be one of: typing, upload_photo, record_video, upload_video, record_audio, upload_audio, upload_document, or find_location
//
// https://core.telegram.org/bots/api#sendchataction
func (b *Bot) SendChatAction(chatId interface{}, action ChatAction) (result ApiResult) {
	// essential params
	params := map[string]interface{}{
		"chat_id": chatId,
		"action":  action,
	}

	return b.requestResult("sendChatAction", params)
}

// Get user profile photos.
//
// options include: offset and limit.
//
// https://core.telegram.org/bots/api#getuserprofilephotos
func (b *Bot) GetUserProfilePhotos(userId int, options map[string]interface{}) (result ApiResultUserProfilePhotos) {
	// essential params
	params := map[string]interface{}{
		"user_id": userId,
	}
	// optional params
	for key, val := range options {
		if val != nil {
			params[key] = val
		}
	}

	return b.requestResultUserProfilePhotos("getUserProfilePhotos", params)
}

// Get updates.
//
// options include: offset, limit, and timeout.
//
// https://core.telegram.org/bots/api#getupdates
func (b *Bot) GetUpdates(options map[string]interface{}) (result ApiResultUpdates) {
	// optional params
	params := map[string]interface{}{}
	for key, val := range options {
		if val != nil {
			params[key] = val
		}
	}

	return b.requestResultUpdates("getUpdates", params)
}

// Get file info and prepare for download.
//
// https://core.telegram.org/bots/api#getfile
func (b *Bot) GetFile(fileId *string) (result ApiResultFile) {
	// essential params
	params := map[string]interface{}{
		"file_id": *fileId,
	}

	return b.requestResultFile("getFile", params)
}

// Get download link from given File.
func (b *Bot) GetFileUrl(file File) string {
	return fmt.Sprintf("%s%s/%s", FileBaseUrl, b.token, *file.FilePath)
}

// Send answers to an inline query.
//
// results = array of InlineQueryResultArticle, InlineQueryResultPhoto, InlineQueryResultGif, InlineQueryResultMpeg4Gif, or InlineQueryResultVideo.
//
// options include: cache_time, is_personal, and next_offset.
//
// https://core.telegram.org/bots/api#answerinlinequery
func (b *Bot) AnswerInlineQuery(inlineQueryId string, results []interface{}, options map[string]interface{}) (result ApiResult) {
	// essential params
	params := map[string]interface{}{
		"inline_query_id": inlineQueryId,
		"results":         results,
	}
	// optional params
	for key, val := range options {
		if val != nil {
			params[key] = val
		}
	}

	return b.requestResult("answerInlineQuery", params)
}

// Check if given http params contain file or not.
func checkIfFileParamExists(params map[string]interface{}) bool {
	for _, value := range params {
		switch value.(type) {
		case *os.File:
			return true
		}
	}

	return false
}

// Convert given interface to string. (for HTTP params)
func (b *Bot) paramToString(param interface{}) (result string, success bool) {
	switch param.(type) {
	case int:
		if intValue, ok := param.(int); ok {
			return strconv.Itoa(intValue), ok
		} else {
			b.error("parameter '%+v' could not be cast to int value", param)
		}
	case float32:
		if floatValue, ok := param.(float32); ok {
			return fmt.Sprintf("%.8f", floatValue), ok
		} else {
			b.error("parameter '%+v' could not be cast to float32 value", param)
		}
	case bool:
		if boolValue, ok := param.(bool); ok {
			return strconv.FormatBool(boolValue), ok
		} else {
			b.error("parameter '%+v' could not be cast to bool value", param)
		}
	case string:
		if strValue, ok := param.(string); ok {
			return strValue, ok
		} else {
			b.error("parameter '%+v' could not be cast to string value", param)
		}
	case []interface{}:
		if json, err := json.Marshal(param); err == nil {
			return string(json), true
		} else {
			b.error(err.Error())
		}
	case ChatAction:
		if value, ok := param.(ChatAction); ok {
			return string(value), ok
		} else {
			b.error("parameter '%+v' could not be cast to string value", param)
		}
	case ParseMode:
		if value, ok := param.(ParseMode); ok {
			return string(value), ok
		} else {
			b.error("parameter '%+v' could not be cast to string value", param)
		}
	case ReplyKeyboardMarkup, ReplyKeyboardHide, ForceReply, InlineQueryResultArticle, InlineQueryResultPhoto, InlineQueryResultGif, InlineQueryResultMpeg4Gif, InlineQueryResultVideo: // (json) structs
		if json, err := json.Marshal(param); err == nil {
			return string(json), true
		} else {
			b.error(err.Error())
		}
	default:
		b.error("unexpected type: '%+v' (%T)", param, param)
	}

	return "", false
}

// Send request to API server and return the response synchronously.
//
// If *os.File is included in the params, it will be closed automatically.
func (b *Bot) sendRequest(method string, params map[string]interface{}) (resp *http.Response, success bool) {
	client := &http.Client{}
	apiUrl := fmt.Sprintf("%s%s/%s", ApiBaseUrl, b.token, method)

	b.verbose("sending request to api url: %s, params: %#v", apiUrl, params)

	if checkIfFileParamExists(params) {
		// multipart form data

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		for key, value := range params {
			switch value.(type) {
			case *os.File:
				if file, ok := value.(*os.File); ok {
					defer file.Close()

					if part, err := writer.CreateFormFile(key, file.Name()); err == nil {
						if _, err = io.Copy(part, file); err != nil {
							b.error("could now write to multipart: %s", key)
						}
					} else {
						b.error("could not create form file for parameter '%s' (%v)", key, value)
					}
				} else {
					b.error("parameter '%s' (%v) could not be cast to file", key, value)
				}
			default:
				if strValue, ok := b.paramToString(value); ok {
					writer.WriteField(key, strValue)
				}
			}
		}

		if err := writer.Close(); err != nil {
			b.error("error while closing writer (%s)", err.Error())
		}

		if req, err := http.NewRequest("POST", apiUrl, body); err == nil {
			req.Header.Add("Content-Type", writer.FormDataContentType()) // due to file parameter

			if resp, err := client.Do(req); err == nil {
				return resp, true
			} else {
				b.error("request error: %s", err.Error())
			}
		} else {
			b.error("building request error: %s", err.Error())
		}
	} else {
		// www-form urlencoded

		paramValues := url.Values{}
		for key, value := range params {
			if strValue, ok := b.paramToString(value); ok {
				paramValues[key] = []string{strValue}
			}
		}
		encoded := paramValues.Encode()

		if req, err := http.NewRequest("POST", apiUrl, bytes.NewBufferString(encoded)); err == nil {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

			if resp, err := client.Do(req); err == nil {
				return resp, true
			} else {
				b.error("request error: %s", err.Error())
			}
		} else {
			b.error("building request error: %s", err.Error())
		}
	}

	return nil, false
}

// Send request for ApiResult and fetch its result.
func (b *Bot) requestResult(method string, params map[string]interface{}) (result ApiResult) {
	var errStr string

	if resp, success := b.sendRequest(method, params); success {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var jsonResponse ApiResult
			if err := json.Unmarshal(body, &jsonResponse); err == nil {
				return jsonResponse
			} else {
				errStr = fmt.Sprintf("json parse error: %s (%s)", err.Error(), string(body))
			}
		} else {
			errStr = fmt.Sprintf("response read error: %s", err.Error())
		}
	} else {
		errStr = fmt.Sprintf("%s failed", method)
	}

	b.error(errStr)

	return ApiResult{Ok: false, Description: &errStr}
}

// Send request for ApiResultUser and fetch its result.
func (b *Bot) requestResultUser(method string, params map[string]interface{}) (result ApiResultUser) {
	var errStr string

	if resp, success := b.sendRequest(method, params); success {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var jsonResponse ApiResultUser
			if err := json.Unmarshal(body, &jsonResponse); err == nil {
				return jsonResponse
			} else {
				errStr = fmt.Sprintf("json parse error: %s (%s)", err.Error(), string(body))
			}
		} else {
			errStr = fmt.Sprintf("response read error: %s", err.Error())
		}
	} else {
		errStr = fmt.Sprintf("%s failed", method)
	}

	b.error(errStr)

	return ApiResultUser{Ok: false, Description: &errStr}
}

// Send request for ApiResultMessage and fetch its result.
func (b *Bot) requestResultMessage(method string, params map[string]interface{}) (result ApiResultMessage) {
	var errStr string

	if resp, success := b.sendRequest(method, params); success {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var jsonResponse ApiResultMessage
			if err := json.Unmarshal(body, &jsonResponse); err == nil {
				return jsonResponse
			} else {
				errStr = fmt.Sprintf("json parse error: %s (%s)", err.Error(), string(body))
			}
		} else {
			errStr = fmt.Sprintf("response read error: %s", err.Error())
		}
	} else {
		errStr = fmt.Sprintf("%s failed", method)
	}

	b.error(errStr)

	return ApiResultMessage{Ok: false, Description: &errStr}
}

// Send request for ApiResultUserProfilePhotos and fetch its result.
func (b *Bot) requestResultUserProfilePhotos(method string, params map[string]interface{}) (result ApiResultUserProfilePhotos) {
	var errStr string

	if resp, success := b.sendRequest(method, params); success {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var jsonResponse ApiResultUserProfilePhotos
			if err := json.Unmarshal(body, &jsonResponse); err == nil {
				return jsonResponse
			} else {
				errStr = fmt.Sprintf("json parse error: %s (%s)", err.Error(), string(body))
			}
		} else {
			errStr = fmt.Sprintf("response read error: %s", err.Error())
		}
	} else {
		errStr = fmt.Sprintf("%s failed", method)
	}

	b.error(errStr)

	return ApiResultUserProfilePhotos{Ok: false, Description: &errStr}
}

// Send request for ApiResultUpdates and fetch its result.
func (b *Bot) requestResultUpdates(method string, params map[string]interface{}) (result ApiResultUpdates) {
	var errStr string

	if resp, success := b.sendRequest(method, params); success {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var jsonResponse ApiResultUpdates
			if err := json.Unmarshal(body, &jsonResponse); err == nil {
				return jsonResponse
			} else {
				errStr = fmt.Sprintf("json parse error: %s (%s)", err.Error(), string(body))
			}
		} else {
			errStr = fmt.Sprintf("response read error: %s", err.Error())
		}
	} else {
		errStr = fmt.Sprintf("%s failed", method)
	}

	b.error(errStr)

	return ApiResultUpdates{Ok: false, Description: &errStr}
}

// Send request for ApiResultFile and fetch its result.
func (b *Bot) requestResultFile(method string, params map[string]interface{}) (result ApiResultFile) {
	var errStr string

	if resp, success := b.sendRequest(method, params); success {
		defer resp.Body.Close()

		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var jsonResponse ApiResultFile
			if err := json.Unmarshal(body, &jsonResponse); err == nil {
				return jsonResponse
			} else {
				errStr = fmt.Sprintf("json parse error: %s (%s)", err.Error(), string(body))
			}
		} else {
			errStr = fmt.Sprintf("response read error: %s", err.Error())
		}
	} else {
		errStr = fmt.Sprintf("%s failed", method)
	}

	b.error(errStr)

	return ApiResultFile{Ok: false, Description: &errStr}
}

// Handle Webhook request.
func (b *Bot) handleWebhook(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	b.verbose("received webhook request: %+v", req)

	if body, err := ioutil.ReadAll(req.Body); err == nil {
		var webhook Update
		if err := json.Unmarshal(body, &webhook); err != nil {
			b.error("error while parsing json (%s)", err.Error())
		} else {
			b.verbose("received webhook body: %s", string(body))

			b.updateHandler(b, webhook, nil)
		}
	} else {
		b.error("error while reading webhook request (%s)", err.Error())

		b.updateHandler(b, Update{}, err)
	}
}
