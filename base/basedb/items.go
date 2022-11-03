package basedb

/***** metadata:VersionObject *****/

const (
	VersionObjectKeyAppIOS     = "app-ios"
	VersionObjectKeyAppAndroid = "app-android"
	VersionObjectKeyWeb        = "app-web"
)

/*
 * ios example by /pub/vobject/app-ios.json
 * {
 *     "code": 0,
 *     "platform": "ios",
 *     "version": "v1.0.0",
 *     "store": "https://apps.apple.com/xxx",
 *     "download": "https://apps.apple.com/xxx"
 * }
 * android example by /pub/vobject/app-android.json
 * {
 *     "code": 0,
 *     "platform": "android",
 *     "version": "v1.0.0",
 *     "store": "https://sj.qq.com/xxx",
 *     "download": "https://xxx.com/xxx"
 * }
 */
type VersionObjectApp struct {
	Platfrom string `json:"platform"` //app platform in ios/android
	Version  string `json:"version"`  //app version
	Store    string `json:"store"`    //app store url address
	Download string `json:"download"` //app download address
	Desc     string `json:"desc"`     //app upgrade description
}
