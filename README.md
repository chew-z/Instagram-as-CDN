# IG CDN

Simple tool for generating list of my Instagram posts using Instagram Basic API from Facebook as [goinsta](https://github.com/ahmdrz/goinsta) that I have been using for a few years [is now dead](https://developers.facebook.com/blog/post/2020/03/10/final-reminder-Instagram-legacy-api-platform-disabled-mar-31/?_fb_noscript=1).

I am using Instagram media URLs for displaying images on my static webpage, it saves me bandwidth as traffic goes directly from Instagram CDN. It is free and legal for personal webpage `[*]` `[**]`.

~~~
To use build main.go and run it. 

Set enviroment variable $IG_TOKEN to your long living token that [you can quickly generate from your Facebook developer account](https://developers.facebook.com/docs/instagram-basic-display-api/overview/#user-token-generator). 

Then go to localhost:8080/json. Parameter limit (default 20) sets how many Instagram post you will get as json. Or use curl (curl 'http://localhost:808/json?limit=10' -o instagram.json' 

~~~

`[*]` _[If you're only using the Legacy API to generate tokens in order to display your Instagram Feed on a personal website, you can use the Instagram Basic Display API's User Token Generator to continue generating tokens. The User Token Generator does not require App Review, but it will only work for your owned Instagram accounts (i.e. accounts for which you know the login credentials](https://www.instagram.com/developer/)_

`[**]` [Instagram Basic Display API](https://developers.facebook.com/docs/instagram-basic-display-api?_fb_noscript=1)


