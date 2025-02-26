package main

import (
	"context"
	"flag"
	"log"
	"math"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var captureURL = flag.String("url", `https://news.qq.com/`, "default is `https://news.qq.com/`")

func main() {
	flag.Parse()

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck,                        //不检查默认浏览器
		chromedp.Flag("headless", true),                       // 不弹出UI界面
		chromedp.Flag("blink-settings", "imagesEnabled=true"), // 加载图片，保证整个完成页面的展示
		chromedp.Flag("ignore-certificate-errors", true),      //忽略错误
		chromedp.Flag("disable-web-security", true),           //禁用网络安全标志
		chromedp.Flag("disable-extensions", true),             //关闭插件支持
		chromedp.Flag("disable-default-apps", true),
		chromedp.WindowSize(1080, 1920),    // 设置浏览器分辨率（窗口大小）--手机窗口尺寸
		chromedp.Flag("disable-gpu", true), // 禁用 GPU加速功能。
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.NoFirstRun, //设置网站不是首次运行
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"), //设置UserAgent
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 创建上下文实例
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// 创建超时上下文
	ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	//导航到目标页面，等待一个元素，捕捉元素的截图
	var buf []byte
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(*captureURL, 100, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("./Screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("图片写入完成")

}

// 获取整个浏览器窗口的截图（全屏）
// 这将模拟浏览器操作设置。
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		//chromedp.WaitVisible("style"),
		chromedp.Sleep(10 * time.Second),
		//chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 得到布局页面
			_, _, _, _, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// 浏览器视窗设置模拟
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// 捕捉屏幕截图
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
