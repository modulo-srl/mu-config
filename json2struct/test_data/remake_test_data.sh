json2struct -declare -pkg test -struct testData settings.default.jsonc > settings.default.def.txt
json2struct -values  -pkg test -struct testData -func initTestData settings.default.jsonc > settings.default.init.txt
