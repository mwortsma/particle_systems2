import matplotlib.pyplot as plt
from matplotlib import pylab
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation
import numpy as np
import sys

def get_keys(distributions):
	for d1 in distributions:
		for d2 in distributions:
			for k in d1:
				if k not in d2: d2[k] = 0
	return sorted(distributions[0].keys())

def plot_path(distributions, labels, show, save, title):
	keys = get_keys(distributions)
	for i in range(len(distributions)):
		d = distributions[i]
		plt.plot([d[k] for k in keys], label=labels[i])
		print [d[k] for k in keys]
	for i in range(len(distributions)):
		for j in range(i+1,len(distributions)):
			print "Error between " + str(labels[i]) + " " + str(labels[j]) + ":"
			d1 = distributions[i]
			d2 = distributions[j]
			print sum([abs(d1[k]-d2[k]) for k in keys])
			print max([abs(d1[k]-d2[k]) for k in keys])
	plt.legend(loc=2)
	plt.xlabel("Possible Neighborhood States")
	plt.ylabel("Probability")
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_time(distributions, labels, show, save, title):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(0,1):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), arr[:,j],
				label=(labels[i]))
	plt.legend(loc=1)
	plt.xlabel("Time")
	plt.ylabel("Probability that a Typical Partilce is Susceptible")
	plt.ylim((0,1))
	plt.title(title)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


def plot_time_full(distributions, labels, show, save, title):
	for i in range(len(distributions)):
		d = distributions[i]
		k = d['K']
		arr = np.array(d['Distr'])
		for j in range(k):
			plt.plot(np.arange(0, len(arr[:,j])*d['Dt'],d['Dt']), arr[:,j],
				label=(labels[i]+" P(X="+str(j))+")")

	plt.legend(loc=2)
	if show:
		plt.show()
	if save and save != "":
		plt.savefig(save)


## GIFFS
def update_time_gif(distributions, labels, ax, total_iters,i):
    ax.cla()
    ax.set_xlabel("Time")
    ax.set_ylabel("Probability")
    arr = np.array(distributions[0]['Distr'])
    T = len(arr[:,0])
    ax.set_xlim((0,T))
    ax.set_ylim((0,1))

    for j in range(0, len(distributions)):
        d = distributions[j]
        k = d['K']
        arr = np.array(d['Distr'])
        T = int((i/float(total_iters))*len(arr[:,0]))
        print T
        ax.plot(np.arange(0, len(arr[:T,0])*d['Dt'],d['Dt']), arr[:T,0],
	        label=(labels[j]+" P(X="+str(0))+")")

	ax.legend(loc=2)
    return None, ax


def plot_time_gif(distributions, labels, show, save, title):
	total_iters = 40

	fig, ax = plt.subplots()
	fig.set_tight_layout(True)

	update = lambda i : update_time_gif(distributions,labels,ax,total_iters,i)

	anim = FuncAnimation(fig, update, frames=np.arange(0, total_iters), interval=200)
	#plt.show()
	anim.save('anim.gif', dpi=80, writer='imagemagick')

def plot_realization(distributions, labels, show, save, title):
	d = np.array(distributions[0]).T

	ax1 = plt.subplot(211)
	ax1.plot(sum(d==0), label='Susceptible', color='red')
	ax1.plot(sum(d==1), label='Infected', color='blue')
	ax1.plot(sum(d==2), label='Recovered', color='green')
	ax1.set_xlabel('Time')
	ax1.set_ylabel('Number of Particles')
	ax1.legend()
	ax1.set_title('SIR Process with a Complete Interaction Network')


	ax2 = plt.subplot(212, sharex=ax1)
	ax2.matshow(1-d, cmap=plt.cm.ocean)
	plt.setp(ax2.get_xticklabels(), visible=False)
	ax2.set_ylabel('Particle')
	ax2.set

	plt.show()
