/**
 * @Time: 2023/10/20 16:48
 * @Author: jzechen
 * @File: contants.go
 * @Software: GoLand collector
 */

package contants

const (
	DefaultK8sClientQPS   = 50
	DefaultK8sClientBurst = 100

	DefaultServerAddr = "0.0.0.0"
	DefaultServerPort = 9302

	DefaultDriveType = "chrome-remote"
	DefaultDriveAddr = "ws://127.0.0.1:9222/devtools/browser"
)
