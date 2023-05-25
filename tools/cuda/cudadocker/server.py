import os
import sys
import json
import time
import importlib

from flask import Flask, request
from fabric import Connection
from invoke.exceptions import UnexpectedExit

app = Flask(__name__)
job_id = None

host = 'pilogin.hpc.sjtu.edu.cn'
user = 'stu1649'
password = 'v8saM*L0'
retry_times = 3

@app.route('/hello', methods=['POST'])
def hello():
    status = dict()
    status['status'] = 'Success'
    return json.dumps(status), 200
    
# 上传并执行任务
@app.route('/sbatch', methods=['POST'])
def sbatch():
    # 获取参数
    config = json.loads(json.dumps(request.json))
    module_name = config['module_name']
    status = dict()
    # 定义变量
    global job_id
    hpc_dir = '/lustre/home/acct-stu/stu1649/gpu'
    cuda_path = f'./tryData/{module_name}.cu'
    slurm_path = f'./tryData/{module_name}.slurm'

    if job_id is None:
        with Connection(host=host, user=user, connect_kwargs={'password' : password}) as c :
            # 上传文件
            for i in range(retry_times):
                try:
                    c.put(cuda_path, remote=hpc_dir)
                    c.put(slurm_path, remote=hpc_dir)
                except Exception as e:
                    print(f'round {i}: {e}')
                else:
                    is_ok = True
                    break
            if not is_ok:
                print('failed to transfer the cuda file')
                sys.exit(1)

            is_ok = False
            for i in range(retry_times):
                try:
                    # 执行文件
                    with c.cd("./gpu"):
                         result = c.run(f'sbatch {module_name}.slurm')
                except UnexpectedExit as ue:
                    print(f'round {i}: {ue}')
                else:
                    job_id = int(result.stdout.strip().split(' ')[-1])
                    is_ok = True
                    break
            if not is_ok:
                print('failed to sbatch slurm script')
                sys.exit(1)

            status['status'] = 'Success'
            status['job_id'] = job_id
    else:
        status['status'] = 'Success'
        status['job_id'] = job_id
    return json.dumps(status), 200

#查询任务状态
@app.route('/query', methods=['POST'])
def query():
    status = dict()
    # 定义变量
    global job_id
    hpc_dir = '/lustre/home/acct-stu/stu1649/gpu'
    is_ok = False
    # 判断查询前有无job被提交过
    if job_id is None:
        print("job_id is none")
        status['status'] = 'Fail'
        status['info'] = 'no job was sbatched before'
    else:
        with Connection(host=host, user=user, connect_kwargs={'password': password}) as c:
            for i in range(retry_times):
                try:
                    c.get(f'{hpc_dir}/out/{job_id}.err', f'./tryData/{job_id}.err')
                    c.get(f'{hpc_dir}/out/{job_id}.out', f'./tryData/{job_id}.out')
                except FileNotFoundError as fe:
                    print(f'round {i}: {fe}')
                    if (f'{fe}' == '[Errno 2] No such file'):
                        status['status'] = 'Running'
                        status['info'] = 'job still running'
                        print("1")
                        return json.dumps(status), 200
                else:
                    is_ok = True
                    break

            if not is_ok:
                status['status'] = 'Error'
                status['info'] = 'something wrong'
                print("2")
                return json.dumps(status), 200

            if os.path.getsize(f'./tryData/{job_id}.err') != 0:
                status['status'] = 'Error'
                with open(f'./tryData/{job_id}.err') as f:
                    status['info'] = f.read()
            else:
                status['status'] = 'Success'
                with open(f'./tryData/{job_id}.out') as f:
                    status['info'] = f.read()
    print("3")
    return json.dumps(status), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8090, processes=True)
