#!/usr/bin/env bash

# TODO: Fix tmux naming

session="ctng"
sleep_duration="0"

# Create new tmux session named ctng
tmux new -d -s $session

# Create CA window
tmux rename-window ca
tmux split-window -h
tmux split-window -h
tmux select-layout tiled

# Run an instance of a CA in each pane
tmux selectp -t 1
tmux send-keys 'echo 1 | ./ca_test.sh 1' C-m
tmux selectp -t 2
tmux send-keys 'echo 1 | ./ca_test.sh 2' C-m
tmux selectp -t 3
sleep $sleep_duration
tmux send-keys 'echo 1 | ./ca_test.sh 3' C-m

# Create logger window
tmux new-window -t 2 -n logger
tmux split-window -h
tmux split-window -h
tmux select-layout tiled

# Run a logger in each pane
tmux selectp -t 1
tmux send-keys 'echo 1 | ./logger_test.sh 1' C-m
tmux selectp -t 2
tmux send-keys 'echo 1 | ./logger_test.sh 2' C-m
tmux selectp -t 3
sleep $sleep_duration
tmux send-keys 'echo 1 | ./logger_test.sh 3' C-m

# Create monitor window
tmux new-window -t 3 -n monitor
tmux split-window -h
tmux split-window -h
tmux split-window -h
tmux select-layout tiled

# Run a monitor in each pane
tmux selectp -t 1
tmux send-keys './monitor_test.sh 1' C-m
tmux selectp -t 2
sleep $sleep_duration
tmux send-keys './monitor_test.sh 2' C-m
tmux selectp -t 3
sleep $sleep_duration
tmux send-keys './monitor_test.sh 3' C-m
tmux selectp -t 4
sleep $sleep_duration
tmux send-keys './monitor_test.sh 4' C-m

# Create gossiper window
tmux new-window -t 4 -n gossiper
tmux split-window -h
tmux split-window -h
tmux split-window -h
tmux select-layout tiled

# Run a gossiper in each pane
tmux selectp -t 1
tmux send-keys './gossiper_test.sh 1' C-m
tmux selectp -t 2
tmux send-keys './gossiper_test.sh 2' C-m
tmux selectp -t 3
sleep $sleep_duration
tmux send-keys './gossiper_test.sh 3' C-m
sleep $sleep_duration
tmux selectp -t 4
tmux send-keys './gossiper_test.sh 4' C-m

# Attach to tmux session
tmux attach-session -t $session