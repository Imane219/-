import subprocess
import os
import re
import shutil
import time
import logging
import sys

sol_path = "./uploads/"
sfuzzdir_path = "./sfuzzdir/"
oyente_path = "./oyente-new/oyente/"
sfuzz_path = "./sFuzz-new/build/fuzzer/"
oyente_output_path = "./oyenteoutput/"
sfuzz_output_path = "./sfuzzoutput/"
max_version = {"0.8": "0.8.1", "0.7": "0.7.6",
            "0.6": "0.6.12", "0.5": "0.5.17", "0.4": "0.4.26"}
sfuzz_test_time = 120

log_file = 'testlog.log'
file_handler = logging.FileHandler(filename=log_file, mode='a')
console_handler = logging.StreamHandler()
file_handler.setLevel('DEBUG')
console_handler.setLevel('DEBUG')
file_handler.setFormatter(logging.Formatter('[%(asctime)s] %(message)s'))
console_handler.setFormatter(logging.Formatter('[%(levelname)s %(funcName)s] %(message)s'))
logger = logging.getLogger()
logger.setLevel('DEBUG')     #设置了这个才会把debug以上的输出到控制台
logger.addHandler(file_handler)
logger.addHandler(console_handler)

def copy_dir(src_dir_path, dest_dir_path):
    if os.path.exists(dest_dir_path):
        shutil.rmtree(dest_dir_path)
    shutil.copytree(src_dir_path, dest_dir_path)


def make_dir(dir_path):
    if os.path.exists(dir_path):
        shutil.rmtree(dir_path)
    os.mkdir(dir_path)


def oyente_generate(sol_path, exec_path):
    """
    oyente运行
    (sol文件存放文件夹, oyente所在文件夹)
    """
    files = os.listdir(sol_path)
    for filename in files:
        if filename[-4:] == '.sol':  # end with".sol"
            logging.debug(f"========oyente start=========")
            # solc版本控制，正则匹配只适用于主版本号为0的情况
            with open(f"{sol_path}/{filename}", "r") as versiontest:
                content = versiontest.read()
                getversion = re.search(
                    "pragma solidity [\^\~](0\.[0-9])([\.0-9]*);\n", content)
                getversion_strict = re.search(
                    "pragma solidity ([\.0-9]*);\n", content)
                if getversion:
                    # 未限死版本
                    # solversion = getversion.group(1)+getversion.group(2)
                    # try to use max version, need test
                    solversion = max_version[getversion.group(1)]
                    subprocess.getstatusoutput(f"solc use {solversion}")
                    logging.debug(f"file {filename}, solc use {solversion}, file version control is {getversion.group(1)+getversion.group(2)}")
                elif getversion_strict:
                    solversion = getversion_strict.group(1)
                    subprocess.getstatusoutput(f"solc use {solversion}")
                    logging.debug(f"file {filename}, solc use {solcversion}, file version control strict")
                else:
                    subprocess.getstatusoutput(f"solc use 0.4.26")
                    logging.debug(f"file {filename}, solc use {solcversion}, no file version control")
            # 编译
            solcbin = subprocess.getstatusoutput(f"solc --bin-runtime --abi '{sol_path}/{filename}' -o ./")
            if (solcbin[0] != 0):
                logging.debug(f"solc return {solcbin[0]}, output:")
                logging.debug(solcbin[1])

            solname = filename[:-4]
            # oyente
            subprocess.getstatusoutput(f"solc use 0.4.22")
            k = subprocess.getstatusoutput(f"python2.7 {exec_path}/oyente.py -s {solname}.bin-runtime -b -ce -j")
            logging.debug(f"oyente return {k[0]}, output:")
            logging.debug(k[1])
            logging.debug(f"========oyente end=========")

    # 删除编译产生的多余文件（需要更新，可能生成多个文件）
    files = os.listdir("./")
    for filename in files:
        if filename[-12:] == '.bin-runtime' or filename[-4:] == ".abi":
            os.remove(filename)


def mk_datdir(sol_path, sfuzzdir_path):
    """
    将oy输出拷贝到sfuzz文件夹
    (sol文件所在路径，sfuzz所在路径)
    """
    logging.debug("copying dir to sfuzz dir...")
    make_dir(sfuzzdir_path)
    dirs = os.listdir("./")
    for solname in dirs:
        if(os.path.isdir(f"./{solname}") and os.path.exists(f"./{solname}/{solname}.txt")):
            copy_dir(f"./{solname}/", f"{sfuzzdir_path}/{solname}/")
            logging.debug(solname)
            shutil.rmtree(f"./{solname}")


def sfuzz_test(sol_dir_path, exec_path, test_time):
    """
    sfuzz初次运行，生成输出(sol文件所在文件夹，sfuzz文件夹，sfuzz单次测试时间)
    """
    logging.debug("========sfuzz start=========")
    make_dir(f"{exec_path}/predicates/")
    make_dir(f"{exec_path}/logger/")
    make_dir(f"{exec_path}/testcases/")
    subprocess.getstatusoutput(f"solc use 0.4.24")

    dirs = os.listdir(sol_dir_path)
    if not dirs:
        logging.debug('No file!')
        return

    sfuzz_sol_dir_path = f"{exec_path}/contracts/"
    pyfile_path = os.getcwd()
    for sol_name in dirs:
        if sol_name.endswith(".sol"):
            os.chdir(pyfile_path)
            make_dir(sfuzz_sol_dir_path)
            shutil.copyfile(f"{sol_dir_path}/{sol_name}",f"{sfuzz_sol_dir_path}/{sol_name}")
            os.chdir(exec_path)
            logging.debug(f"Start fuzzing {sol_name}")
            is_tested = True
            try:
                tmpinfo = subprocess.getstatusoutput(f"./fuzzer -g -r 0 -d {int(test_time)} && chmod +x ./fuzzMe")
                #logging.debug(tmpinfo)
                test_info = subprocess.getstatusoutput(f"./fuzzMe")
            except:
                is_tested = False
            else:
                runtime_info = re.findall(r' run time : (\d+) days, (\d+) hrs, (\d+) min, (\d+) sec', test_info[1])
                if runtime_info:
                    runtime = runtime_info[-1][0] + 'd ' + runtime_info[-1][1] + 'h ' + runtime_info[-1][2] + 'm ' \
                        + runtime_info[-1][3] + 's'
                    logging.debug(f"runtime: {runtime}")

                    coverage_info = re.findall(r'coverage : (\d+%)', test_info[1])[-1]
                    logging.debug(f"cover: {coverage_info}")
                else:
                    is_tested = False

            if not is_tested:
                logging.debug('Not tested')
            
    logging.debug("========sfuzz end=========")


def copy_predicates(sfuzz_path, sol_path):
    """
    将sfuzz产生的预测复制到sol文件对应的文件夹，并给未产生预测的创建一个空文件夹
    """
    plist = os.listdir(f"{sfuzz_path}/predicates")
    for predicate in plist:
        if (predicate.endswith(".txt")):
            solname = predicate[:-4]
            make_dir(f"./{solname}")
            shutil.copyfile(
                f"{sfuzz_path}/predicates/{predicate}", f"./{solname}/{predicate}")
    slist = os.listdir(f"{sol_path}")
    for sol in slist:
        if sol.endswith(".sol"):
            solname = sol[:-4]
            if (solname + ".txt" not in plist):
                make_dir(f"./{solname}")
                subprocess.getstatusoutput(f"touch ./{solname}/{solname}.txt")


def sfuzz_test_end(sol_path, sol_dir_path, exec_path, test_time, sfuzz_output_path):
    """
    sfuzz第二次运行(sol文件和dat所在文件夹，sfuzz文件夹，sfuzz单次测试时间)
    """
    logging.debug("========sfuzz start=========")
    subprocess.getstatusoutput(f"solc use 0.4.24")
    make_dir(f"{exec_path}/logger/")
    make_dir(f"{exec_path}/testcases/")
    make_dir(sfuzz_output_path)

    dirs = os.listdir(sol_dir_path)
    if not dirs:
        logging.debug('No file!')
        return

    sfuzz_sol_dir_path = f"{exec_path}/contracts/"
    pyfile_path = os.getcwd()
    for sol_name in dirs:
        os.chdir(pyfile_path)
        make_dir(sfuzz_sol_dir_path)
        make_dir(f"{exec_path}/testcases/{sol_name}/")
        copy_dir(f"{sol_dir_path}/{sol_name}",
                 f"{exec_path}/testcases/{sol_name}/")
        shutil.copyfile(f"{sol_path}/{sol_name}.sol",
                        f"{sfuzz_sol_dir_path}/{sol_name}.sol")
        os.chdir(exec_path)
        logging.debug('Start fuzzing ' + sol_name)
        is_tested = True
        try:
            tmpinfo = subprocess.getstatusoutput(f"./fuzzer -g -r 0 -d {int(test_time)} && chmod +x ./fuzzMe")
            #logging.debug(tmpinfo)
            test_info = subprocess.getstatusoutput(f'./fuzzMe')
        except:
            is_tested = False
        else:
            runtime_info = re.findall(r' run time : (\d+) days, (\d+) hrs, (\d+) min, (\d+) sec', test_info[1])
            if runtime_info:
                runtime = runtime_info[-1][0] + 'd ' + runtime_info[-1][1] + 'h ' + runtime_info[-1][2] + 'm ' + runtime_info[-1][3] + 's'
                logging.debug(f"runtime: {runtime}")

                coverage_info = re.findall(r'coverage : (\d+%)', test_info[1])[-1]
                logging.debug(f'cover: {coverage_info}')

            else:
                is_tested = False

        if not is_tested:
            logging.debug('Not tested')

        os.chdir(pyfile_path)
        flist = os.listdir(f"{sfuzz_sol_dir_path}")
        for txtoutput in flist:
            if(txtoutput.endswith(r".txt")):
                shutil.copy(f"{sfuzz_sol_dir_path}/{txtoutput}",f"{sfuzz_output_path}/{sol_name}.sol.txt")

    logging.debug("========sfuzz end=========")


def move_oyente_output(oyente_output_path):
    olist = os.listdir("./")
    make_dir(oyente_output_path)
    for jsonoutput in olist:
        if(jsonoutput.endswith(r".bin-runtime.json")):
            shutil.move(f"./{jsonoutput}",f"{oyente_output_path}/{jsonoutput[:-17]}.sol:{jsonoutput[:-17]}.json")

if __name__ == '__main__':
    if len(sys.argv) > 2:
        id = sys.argv[1]
        sol_path += id+"/"
        oyente_output_path += id+"/"
        sfuzz_output_path += id+"/"
        # 前置操作
        fuzz_time = int(sys.argv[2])
        sfuzz_test_time = fuzz_time/2
        pyfile_path = os.getcwd()
        logging.debug(f"\ntest path is {pyfile_path}\ntotal fuzzing time {fuzz_time}\ntesting sol file in {sol_path}...")
        # sfuzz运行
        os.chdir(pyfile_path)
        sfuzz_test(sol_path, sfuzz_path, sfuzz_test_time)
        # 输出预测复制到sol文件所在文件夹
        os.chdir(pyfile_path)
        copy_predicates(sfuzz_path, sol_path)
        # oyente运行（需要-j参数），输出sfuzz需要的分支预测文件
        oyente_generate(sol_path, oyente_path)
        os.chdir(pyfile_path)
        move_oyente_output(oyente_output_path)
        # 复制分支预测到sfuzz位置
        mk_datdir(sol_path, sfuzzdir_path)
        # sfuzz再次运行得到最终输出
        sfuzz_test_end(sol_path, sfuzzdir_path, sfuzz_path, sfuzz_test_time, sfuzz_output_path)
