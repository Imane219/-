import sys
import time

sfuzz_test_time = 4
sol_path = "../uploads/"

def main():
    
    start_time = time.time()
    print("start")
    while True:
        end_time = time.time()
        if end_time-start_time >= sfuzz_test_time:
            break
        print("eafa")
    print("end")


if __name__ == '__main__':
    id = sys.argv[1]
    sol_path += id
    print(sol_path)
    if len(sys.argv) > 2:
        fuzz_time = int(sys.argv[2])
        sfuzz_test_time=fuzz_time/2
    main()