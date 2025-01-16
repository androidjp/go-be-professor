namespace go hello.world

service HelloService {
    string Hello(1: string name) (api.get="/hello")
    string Hello2(1: string name) (api.get="/hello2")
}