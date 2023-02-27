// Code generated by hertz generator.

package Api

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/api/config"
	"github.com/CyanAsterisk/TikGok/server/shared/middleware"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/limiter"
)

func rootMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		// use gzip mw
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
		// use limiter mw
		limiter.AdaptiveLimit(limiter.WithCPUThreshold(900)),
		middleware.Recovery(),
	}
}

func _douyinMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _feedMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getuserinfoMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JWTAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
	}
}

func _publishMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JWTAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
	}
}

func _publishvideoMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _videolistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JWTAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
	}
}

func _comment0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favoriteMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JWTAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
	}
}

func _favorite0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favoritelistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _relationMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JWTAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
	}
}

func __ctionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followinglistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _followerlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _friendMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _friendlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _messageMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		middleware.JWTAuth(config.GlobalServerConfig.JWTInfo.SigningKey),
	}
}

func _chathistoryMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _sentmessageMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _actionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _listMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _action0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _feed0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _action1Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _chatMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _action2Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list1Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _action3Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list2Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list3Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list4Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _login0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _register0Mw() []app.HandlerFunc {
	// your code...
	return nil
}
