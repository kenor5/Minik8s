from flask import Flask, request
import importlib
import json

app = Flask(__name__)

@app.route('/function/<string:function_name>', methods=['POST'])
def execute_function(function_name: str):
    try:
        module = importlib.import_module(function_name)  # 动态导入当前目录下的名字为module_name的模块
        event = {"method": "http"}  # 设置触发器参数，我们当前默认都是http触发
        context = request.get_json()  # 提取POST中的json形式的参数
        print(request)
        function = getattr(module, function_name)  # 获取模块中的函数对象
        result = function(event, context)  # 调用函数并传入参数
        return json.dumps(result), 200  # 将结果转换为JSON字符串，并作为HTTP响应返回
    except Exception as e:
        error_message = {"error": str(e)}
        return json.dumps(error_message), 500  # 返回错误信息并设置HTTP状态码为500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8070, threaded=True)