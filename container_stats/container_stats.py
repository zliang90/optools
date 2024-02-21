#!/usr/bin/env python3

'''
  Container cpu / memory resource usage statistics for docker or containerd.

'''

import subprocess
import json
import os
import sys
from abc import ABC, abstractmethod
from prettytable import PrettyTable


class CRIStats(ABC):
    @abstractmethod
    def sort_by(self):
        pass

    @abstractmethod
    def stats(self):
        pass

    @abstractmethod
    def print_table(self):
        pass


class ContainerStat(CRIStats, object):

    sort_keys = None

    data_rows = None

    limits = 30


    def __init__(self, keys):
        self.sort_by(keys)

        t = PrettyTable()
        t.align = 'l'
        t.field_names = self.field_names
        self.t = t

        self.stats()

    def set_limits(self, n):
        self.limits = n

    def sort_by(self, keys):
        if not keys:
            raise Exception('sort keys is required')

        keys_list = []

        if isinstance(keys, str):
            if ',' in keys:
                keys_list = [ x.strip() for x in keys.strip().split(',') ]
            else:
                keys_list.append(keys.strip())
        elif isinstance(keys, tuple) or isinstance(keys, list):
            for i in keys:
                if isinstance(i, str):
                    keys_list.append(i)

        if len(keys_list) == 0:
            raise Exception('invalid sort by keys')

        sort_keys = []

        for k in keys:
            index = self.index_of_field(k)
            if index == -1:
                raise Exception('invalid sort key: ' + k)

            sort_keys.append(index)

        self.sort_keys = sort_keys

    def index_of_field(self, key):
        for i, v in enumerate(self.field_names):
            if key == v.split('(')[0]:
                return i

        return -1

    def print_table(self):
        if self.data_rows is None:
            raise Exception('no data rows')

        rows = sorted(self.data_rows, key=lambda x : tuple(x[k] for k in self.sort_keys), reverse=True)
        if self.limits > 0:
            self.t.add_rows(rows[:self.limits])
        else:
            self.t.add_rows(rows)

        print(self.t)


class Containerd(ContainerStat):
    ''' containerd stats
    '''

    cmd = ['/usr/bin/crictl', 'stats', '-o', 'json']

    field_names = ['ContainerId',
                   'ContainerName',
                   'PodName',
                   'Namespace',
                   'UsageCoreNanoSeconds',
                   'UsageNanoCores',
                   'CPU(%)',
                   'WorkingSetBytes(MiB)', 'UsageBytes(MiB)', 'RssBytes(MiB)']


    def __init__(self, keys):
        super(Containerd, self).__init__(keys)

    def __format_int(self, x, y):
        if not isinstance(x, dict) or not isinstance(y, str):
            raise Exception('invalid type, x dict , y str expected')

        return int(x[y]['value'])

    def __format_bytes(self, x, y):
        if not isinstance(x, dict) or not isinstance(y, str):
            raise Exception('invalid type, x dict , y str expected')

        return float('{:.2f}'.format(self.__format_int(x, y)/1000/1000))

    def stats(self):
        r = subprocess.run(self.cmd, capture_output=True)

        if r.returncode != 0:
            raise Exception(r.stderr.decode())

        rows = []

        try:
            output = r.stdout.decode().strip()
            data = json.loads(output)

            for c in data['stats']:
                cpu, memory = c['cpu'], c['memory']
                if cpu is None and memory is None:
                    continue

                container_id, name = c['attributes']['id'], c['attributes']['metadata']['name']
                labels = c['attributes']['labels']
                pod_name = labels['io.kubernetes.pod.name']
                namespace = labels['io.kubernetes.pod.namespace']

                row = [ None for _ in range(len(self.field_names)) ]
                row[0], row[1] = container_id[:13], name
                row[2], row[3] = pod_name, namespace

                if cpu:
                    row[4], row[5] = self.__format_int(cpu, 'usageCoreNanoSeconds'), self.__format_int(cpu, 'usageNanoCores')
                    row[6] = round((self.__format_int(cpu, 'usageNanoCores')/(10**9))*100, 2)

                if memory:
                    row[7] = self.__format_bytes(memory, 'workingSetBytes')
                    row[8], row[9] = self.__format_bytes(memory, 'usageBytes'), self.__format_bytes(memory, 'rssBytes')
                rows.append(row)

            # table rows
            self.data_rows = rows
        except:
            raise


class Docker(ContainerStat):
    ''' docker stats
    '''

    cmd = ['docker',
           'stats',
           '--no-stream',
           '--format',
           '{"container":"{{ .Container }}","name":"{{ .Name }}","memory":{"raw":"{{ .MemUsage }}","percent":"{{ .MemPerc }}"},"cpu":"{{ .CPUPerc }}"}'
        ]

    field_names = ['ContainerId',
                   'ContainerName',
                   'PodName',
                   'Namespace',
                   'CPU(%)',
                   'Memory(%)',
                   'MemoryUsed(MiB)',
                   'MemoryTotal(MiB)',
                ]


    def __init__(self, keys):
        super(Docker, self).__init__(keys)

    def stats(self):
        r = subprocess.run(self.cmd, capture_output=True)

        if r.returncode != 0:
            raise Exception(r.stderr.decode())

        try:
            rows = []

            output = r.stdout.decode()
            for line in output.split('\n'):
                if not line:
                    continue

                row = [ None for _ in range(len(self.field_names)) ]
                json_row = json.loads(line.strip())

                if 'cpu' not in json_row or 'memory' not in json_row:
                    continue

                row[0], name_fields = json_row['container'], json_row['name'][4:].split('_')

                container_name, pod_name, namespace =  name_fields[0], name_fields[1], name_fields[2]
                if container_name == 'POD':
                    continue

                row[1], row[2], row[3] =  container_name, pod_name, namespace

                cpu_percent = json_row['cpu'].rstrip('%')
                if cpu_percent:
                    row[4] = float(cpu_percent)

                memory = json_row['memory']
                if memory:
                    raw, percent = memory['raw'], memory['percent'].rstrip('%')
                    if percent:
                        row[5] = float(percent)
                    raws = [ x.strip() for x in raw.split('/') ]
                    if len(raws) == 2:
                        row[6], row[7] = self.format_to_mb(raws[0]), self.format_to_mb(raws[1])

                rows.append(row)
            # add data rows
            self.data_rows = rows
        except:
            raise

    def index_alphabet(self, v):
        for i, c in enumerate(v):
            if (ord(c) <= 122 and ord(c) >= 97) or \
                    (ord(c) <= 90 and ord(c) >= 65):
                return (float(v[:i]), v[i:])

        return None

    def format_to_mb(self, v):
        size_unit = self.index_alphabet(v)

        if size_unit is None:
            raise Exception('invalid size or unit of "{}"'.format(v))

        size, unit = size_unit

        if unit == 'KiB':
            return size / 1000
        elif unit == 'MiB':
            return size
        elif unit == 'GiB':
            return size * 1000
        elif unit == 'TiB':
            return size * 1000 * 1000
        else:
            raise Exception('unsported unit "{}"'.foamt(unit))


if __name__ == '__main__':
    c = None
    limits = 30
    sort_keys = []

    if len(sys.argv) >= 2:
        keys = sys.argv[1]
        sort_keys = [ x.strip() for x in keys.split(',') ]
    if len(sys.argv) >= 3:
        limits = int(sys.argv[2])

    if os.path.exists('/run/containerd/containerd.sock'):
        if len(sort_keys) == 0:
            sort_keys = ['WorkingSetBytes', 'CPU']
        c = Containerd(sort_keys)
    elif os.path.exists('/var/run/docker.sock'):
        if len(sort_keys) == 0:
            sort_keys = ['Memory', 'CPU']
        c = Docker(sort_keys)
    else:
        print('no cri docker or containerd')
        sys.exit(1)

    c.set_limits(limits)
    c.print_table()
