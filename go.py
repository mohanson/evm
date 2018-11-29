import subprocess
import sys


def call(command):
    print(command)
    r = subprocess.call(command, shell=True)
    if r != 0:
        sys.exit(r)


def main():
    call(f'go install github.com/mohanson/evmfun/cmd/evmfun')


if __name__ == '__main__':
    main()
