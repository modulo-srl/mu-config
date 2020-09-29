json2struct -declare -pkg test -struct testData settings.default.json > settings.default.def.txt
json2struct -values  -pkg test -struct testData -func initTestData settings.default.json > settings.default.init.txt
