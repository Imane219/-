import subprocess
import os
import re
import shutil
import time
import sys

sol_path="../uploads/"
sfuzzdir_path="../sfuzzdir"
oyente_path="../../oyente-new/oyente"
oyente_output_file_path="./oyenteoutput"
oyente_log_file_path="./oyentelog"
sfuzz_path="../../sFuzz-new/build/fuzzer"
sfuzz_output_file_path = "./sfuzzoutput"
max_version = {"0.8":"0.8.1", "0.7":"0.7.6", "0.6":"0.6.12", "0.5":"0.5.17", "0.4":"0.4.26"}

sfuzz_test_time = 120
sfuzz_bug_list = ['gasless send', 'exception disorder', 'reentrancy', 'timestamp dependency',
              'block number dependency', 'dangerous delegatecall', 'freezing ether',
              'integer overflow', 'integer underflow']

def copy_dir(src_dir_path, dest_dir_path):
    if os.path.exists(dest_dir_path):
        shutil.rmtree(dest_dir_path)
    shutil.copytree(src_dir_path, dest_dir_path)

def make_dir(dir_path):
    if os.path.exists(dir_path):
        shutil.rmtree(dir_path)
    os.mkdir(dir_path)

def oyente_generate(sol_path,exec_path,output_file_path,log_file_path):
    out=open(output_file_path,"w")
    log=open(log_file_path,"w")
    out.write("filename, err code, Callstack Depth Attack Vulnerability, Transaction-Ordering Dependence (TOD), Timestamp Dependency, Re-Entrancy Vulnerability, \n")
    files = os.listdir(sol_path)
    for filename in files:
        if filename[-4:] == '.sol':  # end with".sol"
            """
            # solc版本控制，正则匹配只适用于主版本号为0的情况
            with open(f"{sol_path}/{filename}","r") as versiontest:
                content=versiontest.read()
                getversion = re.search("pragma solidity [\^\~](0\.[0-9])([\.0-9]*);\n", content)
                getversion_strict = re.search("pragma solidity ([\.0-9]*);\n", content)
                if getversion:
                    solversion = getversion.group(1)+getversion.group(2)
                    print(solversion)
                    subprocess.getstatusoutput(f"solc use {solversion}")
                elif getversion_strict:
                    solversion = getversion_strict.group(1)
                    print(solversion)
                    subprocess.getstatusoutput(f"solc use {solversion}")
                else:
                    subprocess.getstatusoutput(f"solc use 0.4.26")
            
            # 编译
            solcbin = subprocess.getstatusoutput(f"solc --bin-runtime --abi '{sol_path}/{filename}' -o ./")
            if (solcbin[0] != 0):
                log.write(f"solc output:{solcbin}\n")
            """
            solname = filename[:-4]
            
            # oyente
            subprocess.getstatusoutput(f"solc use 0.4.22")
            k=subprocess.getstatusoutput(f"python2.7 {exec_path}/oyente.py -s {solname}.bin-runtime -b -ce")
            out.write(f"{filename}, {k[0]}, ")

            # 结果读取
            pattern="INFO:symExec:( |\t)*(.*):( |\t)*(True|False)"
            matches=re.finditer(pattern, k[1], flags=0)
            bugs = [(match.group(2),match.group(4)) for match in matches]
            for bug in bugs:
                out.write(f"{bug[1]}, ")
            out.write("\n")
            print(f"scan {filename}, return {k[0]}")
            log.write((f"scan {filename}, return {k[0]}\n"))
            log.write(f"=================\nlog:\n{k[1]}\n=================\n\n")

    # 删除编译产生的多余文件
    files = os.listdir("./")
    for filename in files:
        if filename[-12:] == '.bin-runtime' or filename[-4:] == ".abi":
            os.remove(filename)
    
    out.close()
    log.close()

def mk_datdir(sol_path,sfuzzdir_path):
    #copy file to sfuzzdir
    make_dir(sfuzzdir_path)
    dirs = os.listdir("./")
    for solname in dirs:
        if os.path.isdir(f"./{solname}"):
            copy_dir(f"./{solname}",f"{sfuzzdir_path}/{solname}")
            print(solname)
            shutil.rmtree(f"./{solname}")

def sfuzz_test(sol_dir_path, exec_path, output_path, test_time, bug_list):
    make_dir(f"{exec_path}/predicates/")
    subprocess.getstatusoutput(f"solc use 0.4.24")
    result_file=open(output_path,"w")
    result_file.write('contract name, Run Time, Coverage, ')
    for bugname in bug_list:
        result_file.write(f"{bugname}, ")
    result_file.write("Not Tested, \n")

    dirs = os.listdir(sol_dir_path)
    if not dirs:
        print('No file!')
        return

    sfuzz_sol_dir_path = f"{exec_path}/contracts/"
    pyfile_path=os.getcwd()
    for sol_name in dirs:
        if sol_name.endswith(".sol"):
            os.chdir(pyfile_path)
            make_dir(sfuzz_sol_dir_path)
            shutil.copyfile(f"{sol_dir_path}/{sol_name}",sfuzz_sol_dir_path+sol_name)
            os.chdir(exec_path)
            print('Start fuzzing ' + sol_name)

            result_file.write(f"{sol_name}, ")
            print(sol_name, end=': ')
            is_tested = True
            try:
                subprocess.getstatusoutput(f'./fuzzer -g -r 0 -d {test_time} && chmod +x ./fuzzMe')
                test_info = subprocess.getstatusoutput(f'./fuzzMe')
            except:
                is_tested = False
            else:
                runtime_info = re.findall(r' run time : (\d+) days, (\d+) hrs, (\d+) min, (\d+) sec', test_info[1])
                if runtime_info:
                    runtime = runtime_info[-1][0] + 'd ' + runtime_info[-1][1] + 'h ' + runtime_info[-1][2] + 'm ' \
                            + runtime_info[-1][3] + 's'
                    result_file.write(f"{runtime}, ")
                    print(runtime, end='  ')

                    coverage_info = re.findall(r'coverage : (\d+%)', test_info[1])[-1]
                    result_file.write(f"{coverage_info}, ")
                    print(f'cover:{coverage_info}')

                    for bug in bug_list:
                        info = re.search(f'{bug} : found', test_info[1])
                        result_file.write(f"{bool(info)}, ")
                        print(f'{bug}:{bool(info)}', end=' ')
                else:
                    is_tested = False

            result_file.write(f"{is_tested}, \n")
            if not is_tested:
                print('Not tested', end='')
            print('\n')
    
    result_file.close()

def copyPredicates(sfuzz_path,sol_path):
    plist = os.listdir(f"{sfuzz_path}/predicates")
    for predicate in plist:
        if (predicate.endswith(".txt")):
            solname = predicate[:-4]
            make_dir(f"./{solname}")
            shutil.copyfile(f"{sfuzz_path}/predicates/{predicate}", f"./{solname}/{predicate}")
    slist = os.listdir(f"{sol_path}")
    for sol in slist:
        if sol.endswith(".sol"):
            solname = sol[:-4]
            if (solname + ".txt" not in plist):
                make_dir(f"./{solname}")
                subprocess.getstatusoutput(f"touch ./{solname}/{solname}.txt")

def sfuzz_test_end(sol_dir_path, exec_path, output_path,test_time,bug_list):
    subprocess.getstatusoutput(f"solc use 0.4.24")
    result_file=open(output_path,"w")
    result_file.write('contract name, Run Time, Coverage, ')
    for bugname in bug_list:
        result_file.write(f"{bugname}, ")
    result_file.write("Not Tested, \n")

    dirs = os.listdir(sol_dir_path)
    if not dirs:
        print('No file!')
        return

    sfuzz_sol_dir_path = f"{exec_path}/contracts/"
    pyfile_path=os.getcwd()
    for sol_name in dirs:
        os.chdir(pyfile_path)
        make_dir(sfuzz_sol_dir_path)
        make_dir(f"{exec_path}/testcases/{sol_name}/")
        copy_dir(f"{sol_dir_path}/{sol_name}", f"{exec_path}/testcases/{sol_name}/")
        shutil.copyfile(f"{sol_path}/{sol_name}.sol",f"{sfuzz_sol_dir_path}/{sol_name}.sol")
        os.chdir(exec_path)
        print('Start fuzzing ' + sol_name)

        result_file.write(f"{sol_name}, ")
        print(sol_name, end=': ')
        is_tested = True
        try:
            subprocess.getstatusoutput(f'./fuzzer -g -r 0 -d {test_time} && chmod +x ./fuzzMe')
            test_info = subprocess.getstatusoutput(f'./fuzzMe')
        except:
            is_tested = False
        else:
            runtime_info = re.findall(r' run time : (\d+) days, (\d+) hrs, (\d+) min, (\d+) sec', test_info[1])
            if runtime_info:
                runtime = runtime_info[-1][0] + 'd ' + runtime_info[-1][1] + 'h ' + runtime_info[-1][2] + 'm ' \
                        + runtime_info[-1][3] + 's'
                result_file.write(f"{runtime}, ")
                print(runtime, end='  ')

                coverage_info = re.findall(r'coverage : (\d+%)', test_info[1])[-1]
                result_file.write(f"{coverage_info}, ")
                print(f'cover:{coverage_info}')

                for bug in bug_list:
                    info = re.search(f'{bug} : found', test_info[1])
                    result_file.write(f"{bool(info)}, ")
                    print(f'{bug}:{bool(info)}', end=' ')
            else:
                is_tested = False

        result_file.write(f"{is_tested}, \n")
        if not is_tested:
            print('Not tested', end='')
        print('\n')
    
    result_file.close()

if __name__ == '__main__':
    id = sys.argv[1]
    sol_path += id+"/"
    oyente_output_file_path += id +"/"
    sfuzz_output_file_path += id + "/"
    if len(sys.argv) > 2:
        fuzz_time = int(sys.argv[2])
        sfuzz_test_time=fuzz_time/2

    pyfile_path=os.getcwd()
    vuldirs = os.listdir("../dataset")
    for vultype in vuldirs:
        os.chdir(pyfile_path)
        sfuzz_test(f"../dataset/{vultype}", sfuzz_path, f"{sfuzz_output_file_path}-{vultype}.csv", sfuzz_test_time, sfuzz_bug_list)
        os.chdir(pyfile_path)
        copyPredicates(sfuzz_path,f"../dataset/{vultype}")
        oyente_generate(f"../dataset/{vultype}", oyente_path, f"{oyente_output_file_path}-{vultype}.csv", f"{oyente_log_file_path}-{vultype}.txt")
        os.chdir(pyfile_path)
        mk_datdir(sol_path,sfuzzdir_path)
        copy_dir(f"../dataset/{vultype}/",f"{sol_path}")
        sfuzz_test_end(sfuzzdir_path,sfuzz_path,f"{sfuzz_output_file_path}_end-{vultype}.csv",sfuzz_test_time,sfuzz_bug_list)
