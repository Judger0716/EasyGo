module EasyGo

go 1.21.5

require (
    EasyCache v0.0.0
    EasyGin v0.0.0
)

replace (
    EasyCache => ./EasyCache
    EasyGin => ./EasyGin
)