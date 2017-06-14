#!/bin/env python

from __future__ import print_function
import os
import sys
import random
import datetime
import time
import xml.etree.ElementTree as ET
from pytz import timezone


SCRIPT_PATH = os.path.split(os.path.realpath(sys.argv[0]))[0]
os.chdir(SCRIPT_PATH)

START_TIME = time.time()


def how_long_has_it_been():
    return int(time.time() - START_TIME) / 60


def retry_rank(team_list):
    sort_list = []
    for team in team_list:
        triple = [
            int(team.get('solved')),
            int(team.get('points')),
            team.get('teamKey'),
        ]
        sort_list.append(triple)

    sort_list = sorted(sort_list,
                       key=lambda triple: triple[0] * 10000 - triple[1], reverse=True)
    rank_table = {}
    for i, triple in enumerate(sort_list):
        rank_table[triple[2]] = i + 1

    for team in team_list:
        rank = rank_table[team.get('teamKey')]
        team.set('index', str(rank))
        team.set('rank', str(rank))


def solve(psi, team, problem_list):
    if psi.get('isSolved') == 'true':
        return

    time_spent = how_long_has_it_been()
    problem = [p for p in problem_list if p.get('id') == psi.get('index')][0]

    attempts = int(psi.get('attempts'))
    if psi.get('isPending') != 'true':
        attempts += 1
        psi.set('attempts', str(attempts))
        problem.set('attempts', str(int(problem.get('attempts')) + 1))
    else:
        psi.set('isPending', 'false')
    psi.set('isSolved', 'true')
    points = time_spent + (attempts - 1) * 20
    psi.set('points', str(points))
    team.set('points', str(int(team.get('points')) + points))
    team.set('solved', str(int(team.get('solved')) + 1))

    numberSolved = int(problem.get('numberSolved'))
    numberSolved += 1
    if numberSolved == 1:
        problem.set('bestSolutionTime', str(time_spent))
    problem.set('lastSolutionTime', str(time_spent))
    problem.set('numberSolved', str(numberSolved))


def pending(psi, team, problem_list):
    if psi.get('isSolved') == 'true':
        return

    time_spent = how_long_has_it_been()
    problem = [p for p in problem_list if p.get('id') == psi.get('index')][0]

    attempts = int(psi.get('attempts'))
    if psi.get('isPending') != 'true':
        attempts += 1
        psi.set('attempts', str(attempts))
        psi.set('isPending', 'true')
        problem.set('attempts', str(int(problem.get('attempts')) + 1))


def wrong(psi, team, problem_list):
    if psi.get('isSolved') == 'true':
        return

    time_spent = how_long_has_it_been()
    problem = [p for p in problem_list if p.get('id') == psi.get('index')][0]

    attempts = int(psi.get('attempts'))
    if psi.get('isPending') != 'true':
        attempts += 1
        psi.set('attempts', str(attempts))
        problem.set('attempts', str(int(problem.get('attempts')) + 1))
    else:
        psi.set('isPending', 'false')


def mock(results_xml):
    print('mock', file=sys.stderr)
    tree = ET.parse(results_xml)
    root = tree.getroot()
    standings_header = tree.find('standingsHeader')
    problem_list = standings_header.findall('problem')
    team_list = root.findall('teamStanding')

    action_num = random.randint(0, 10)
    for _ in range(action_num):
        tid = random.randint(0, len(team_list) - 1)
        team = team_list[tid]
        psi_list = team.findall('problemSummaryInfo')
        pid = random.randint(0, len(psi_list) - 1)
        psi = psi_list[pid]

        dice = random.random()
        if dice < 1.0 / 3.0:
            pending(psi, team, problem_list)
        elif 1.0 / 3.0 <= dice < 2.0 / 3.0:
            solve(psi, team, problem_list)
        else:
            wrong(psi, team, problem_list)

    retry_rank(team_list)
    time_format = '%a %b %d %H:%M:%S %Z %Y'
    last = datetime.datetime.strptime(
        standings_header.get('currentDate'),
        time_format,
    ).replace(tzinfo=timezone('CST6CDT'))
    now = last + datetime.timedelta(seconds=int(time.time() - START_TIME))
    standings_header.set('currentDate', now.strftime(time_format))
    tree.write(results_xml)


def main():
    results_xml = os.path.join(SCRIPT_PATH, 'results.xml')
    while True:
        time.sleep(random.uniform(0.5, 4))
        mock(results_xml)


if __name__ == '__main__':
    main()
