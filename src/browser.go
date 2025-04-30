package main

import (
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
)

func getNewBrowser() *browser.Browser {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Firefox())
	bow.SetAttributes(browser.AttributeMap{
		browser.SendReferer:         false,
		browser.MetaRefreshHandling: true,
		browser.FollowRedirects:     true,
	})
	bow.SetCookieJar(jar.NewMemoryCookies())
	return bow
}
