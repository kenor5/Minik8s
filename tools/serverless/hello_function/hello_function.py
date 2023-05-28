def hello_function(event: dict, context: dict)->dict:
    name = context['name']
    age = context['age']

    print(name)  # 输出：John
    print(age)  # 输出：18

    return {"result": "hello, {} years old guy {}".format(age, name)}