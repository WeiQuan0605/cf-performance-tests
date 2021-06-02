#!/usr/bin/env python

import argparse
import json
from colorama import init, Fore, Style
from tabulate import tabulate


init()

parser = argparse.ArgumentParser(description="Compare CF performance test results")
parser.add_argument(
    "--previous",
    required=True,
    type=open,
    help="previous test results to compare",
)
parser.add_argument(
    "--current",
    required=True,
    type=open,
    help="current test results to compare",
)
parser.add_argument(
    "--threshold",
    default=110,
    type=int,
    help="maximum allowed percentage difference between test results (default 110%)",
)

args = parser.parse_args()

threshold = args.threshold / 100
previous = json.load(args.previous)
current = json.load(args.current)

failures = []
for test, prev_results in previous.items():
    for metric in ["Largest", "Average", "StdDeviation"]:
        prev_metric = prev_results["request time"][metric]
        curr_metric = current[test]["request time"][metric]
        if curr_metric > (threshold * prev_metric):
            failures.append((test, metric, curr_metric / prev_metric * 100))

if len(failures) > 0:
    print(
        f"\n{Fore.RED}{Style.BRIGHT} ERROR: {len(failures)} test metric(s) were not within {args.threshold}% threshold {Style.RESET_ALL}"
    )
    print(
        tabulate(
            failures,
            headers=["Test", "Metric", "Percentage of previous"],
            floatfmt=".2f",
            tablefmt="fancy_grid",
        )
    )
    exit(1)

print(
    f"\n{Fore.GREEN}{Style.BRIGHT} All test metrics were within {args.threshold}% threshold"
)
