import re
import sys
import re
import subprocess
from os import walk
from os.path import join
import os


def cli_run(cmd: str):
    with subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE) as p:
        if p.returncode != None and p.returncode != 0:
            raise ValueError(f"'{cmd}' failed, returncode {p.returncode}")
        std = str(p.stdout.read(), 'utf-8')
        return std


def current_branch():
    out = cli_run("git status")
    lines = out.splitlines()
    for l in lines:
        m = re.match('On branch ([^\\s]+)', l)
        if m:
            return m[1]


def all_tags():
    return cli_run("git tag")


def current_tag():
    out = cli_run("git describe --tags --abbrev=0")
    return out.strip()


def parse_beta(tag):
    pat = re.compile('(v.+).beta.*')
    m = pat.match(tag)
    if m:
        return m[1]
    return tag


def guess_next(last):
    if not last:
        return "v0.0.1"
    if is_beta(last):
        return incr_beta(last)
    return incr_release(last)


def incr_release(v):
    pat = re.compile('(v?\\d+\\.\\d+\\.)(\\d+)\\.?.*')
    m = pat.match(v)
    return m[1] + str(int(m[2]) + 1)


def incr_beta(v):
    pat = re.compile('^v?\\d+\\.\\d+\\.\\d+$')
    m = pat.match(v)
    if m:
        return v + "-beta.1"

    pat = re.compile('(v?\\d+\\.\\d+\\.\\d+-beta\\.)(\\d+)')
    m = pat.match(v)
    return m[1] + str(int(m[2]) + 1)


def is_beta(v):
    pat = re.compile('(v?\\d+\\.\\d+\\.\\d+-beta\\.)(\\d+)')
    return pat.match(v)


if __name__ == '__main__':

    target = None

    for v in sys.argv[1:]:
        if not v:
            continue
        if v == "-f":
            print("-f is no longer supported: this script only updates the version constants. Commit and tag yourself.")
            exit(1)
        target = v

    if not target:
        last = current_tag()
        guessed = guess_next(last)
        print()
        possible_options = []
        possible_options.append("'" + guessed + "'")

        if "beta" not in guessed:
           possible_options.append("'" +  guessed + "-beta.1" + "'")
        else:
           possible_options.append("'" + guessed[:guessed.rfind("-beta")] + "'")

        opts = " or ".join(possible_options)
        print(f"Last release is '{last}'. Are you planning to release {opts} ?")
        print()
        print(f"release {guessed}")
        print()
        exit(1)

    if target:
        latest_tag = current_tag().strip()
        if latest_tag == target:
            print(f"{latest_tag} has been released")
            exit(1)

        target_release = parse_beta(target)
        if target_release != target and target_release in all_tags().splitlines():
            print(f"{target_release} has been released")
            exit(1)

    for (dir_path, dir_name, file_names) in walk("."):
        for fn in file_names:
            if fn == 'version.go':
                version_file = join(dir_path, fn)
                print(version_file)
                pkg = "main"
                matched = False

                pat = re.compile('package (.*)')
                with open(version_file, "r") as f:
                    lines = f.readlines()
                    for l in lines:
                        m = pat.match(l)
                        if m:
                            pkg = m[1]
                        if re.match("\\s*Version *= *\".*\"\\s*", l):
                            matched = True

                # print("matched? ", matched)
                if matched:
                    with open(version_file, "w") as f:
                        f.writelines([
                            f"package {pkg}\n",
                            "\n",
                            "var (\n",
                            f"\tVersion = \"{target}\"\n"
                            ")\n"
                            ""
                        ])

            if fn == 'version.ts':
                version_file = join(dir_path, fn)
                # print(version_file)
                matched = False

                with open(version_file, "r") as f:
                    lines = f.readlines()
                    for l in lines:
                        if re.match("\\s*export +const +Version *= *[\"'].*[\"']\\s*", l):
                            matched = True

                # print("matched? ", matched)
                if matched:
                    with open(version_file, "w") as f:
                        f.writelines([
                            f"export const Version = \"{target}\"\n"
                        ])

    mod_file = None
    for (dir_path, dir_name, file_names) in walk("."):
        for fn in file_names:
            if fn == 'go.mod':
                print(f"\nrunning 'go fmt ./...' in {dir_path}/")
                print(cli_run(f"cd {dir_path} && go fmt ./..."))
                break

    wd = os.getcwd()
    if "beta" not in target and "miso" in wd:
        for (dir_path, dir_name, file_names) in walk("."):
            for fn in file_names:
                if fn in ['README.md', 'tools.md']:
                    with open(join(dir_path, fn), "r+") as f:
                        all = f.read()
                        m = re.compile("(go +(?:get|install) .*)(v\\d+\\.\\d+\\.\\d+(?:-beta\\.\\d+)?)")
                        all = re.sub(m, r"\1" + target, all)
                        f.truncate(0)
                        f.seek(0)
                        f.write(all)
                    break


    print()
    print(f"Version constants updated to {target}. Nothing was committed or tagged.")
    print("Review the diff, then commit and tag yourself:")
    print("  git add -u")
    print(f"  git commit -m \"Release {target}\"")
    print(f"  git tag {target}")
    print(f"  git push && git push origin {target}")
