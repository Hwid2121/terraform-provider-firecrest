#!/bin/bash -l
#SBATCH --job-name="coder-job"
#SBATCH --mail-type=ALL
#SBATCH --mail-user="nicolotafta@gmail.com"
#SBATCH --time=00:08:00
#SBATCH --nodes=1
#SBATCH --ntasks-per-core=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=6
#SBATCH --partition=debug
#SBATCH --constraint=gpu

mkdir -p $SCRATCH/firecrest/$SLURM_JOBID

# Get node IP
node_name=$(scontrol show hostname $SLURM_JOB_NODELIST)
node_ip=$(getent hosts $node_name | awk '{ print $1 }')

# Log the node IP
echo "Node name: $node_name"
echo "Node IP: $node_ip"
echo $node_ip > $SCRATCH/firecrest/$SLURM_JOBID/node_ip.txt
echo "IP address written to $SCRATCH/firecrest/$SLURM_JOBID/node_ip.txt"
    cat - > agent.sh <<< ''
    chmod +x agent.sh

    module load daint-gpu
    module load sarus

    node_name=$(scontrol show hostname $SLURM_JOB_NODELIST)
    node_ip=$(getent hosts $node_name | awk '{ print $1 }')

    echo "Node name: $node_name"
    echo "Node IP: $node_ip"

    export CODER_AGENT_TOKEN=3cb67f28-f5d4-4313-bf4a-835da072eb21
    export CODER_AGENT_ID=640011a9-1516-4722-937f-b3adc004a0c9

    echo "Coder Token: $CODER_AGENT_TOKEN "
    echo "Coder ID: $CODER_AGENT_ID"


    srun sarus pull nikotaft/coder-environment:latest
    srun sarus run nikotaft/coder-environment:latest /bin/bash -c "
      curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/tmp/code-server --edge &&
env SHELL=/bin/bash HOME=/home/coder PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin /tmp/code-server/bin/code-server --auth none --port 8080 --log debug    " &
     
     ./agent.sh
     echo "Agent is not blocking..."
     sleep 3600

