.. _pmm.glossary.pmm-server.additional-option:

================================================================================
|pmm-server| Additional Options 
================================================================================

This glossary contains the addtional parameters that you may pass when starting
|pmm-server|.

If you use |docker| to run the server, use the :option:`-e` flag followed by the
parameter. Use this flag in front of each parameter that you pass to
|pmm-server|.

Here, we pass more than one option to |pmm-server| along with the |docker.run|
command. |tip.run-this.root|.

.. include:: .res/code/docker.run.server-user.example.txt

List of |pmm-server| Parameters
================================================================================

.. glossary::
   :sorted:

   DISABLE_TELEMETRY

      With :term:`telemetry` enabled, your |pmm-server| sends some statistics to
      `v.percona.com`_ every 24 hours. This statistics includes the following
      details:

      - |pmm-server| unique ID
      - |pmm| version
      - The name and version of the operating system
      - |mysql| version
      - |perl| version

      If you do not want your |pmm-server| to send this information, disable telemetry
      when running your |docker| container:

      .. include:: .res/code/docker.run.disable-telemetry.txt

   METRICS_RETENTION

      This option determines how long metrics are stored at :term:`PMM
      Server`. The value is passed as a combination of hours, minutes, and
      seconds, such as **720h0m0s**. The minutes (a number followed by *m*) and
      seconds (a number followed by *s*) are optional.

      To set the |opt.metrics-retention| option to 8 days, set this option to *192h*.

      |tip.run-this.root|

      .. include:: .res/code/docker.run.e.metrics-retention.txt

      .. seealso::

	 Data retention in PMM
	    :term:`Data retention`
	 Queries retention
	    :term:`QUERIES_RETENTION <QUERIES_RETENTION>`

   QUERIES_RETENTION

      This option determines how many days queries are stored at :term:`PMM Server`. 

      .. include:: .res/code/docker.run.e.queries-retention.txt

      .. seealso::

	 Metrics retention
	    :term:`METRICS_RETENTION <METRICS_RETENTION>`
	 Data retention in PMM
	    :term:`Data retention`

   ORCHESTRATOR_ENABLED

      This option enables |orchestrator| (See
      :ref:`pmm.using.orchestrator`). By default it is disabled. It is
      also desabled if this option contains **false**.

      .. include:: .res/code/docker.run.orchestrator-enabled.txt

      .. seealso::

	 Orchestrator
	    :term:`Orchestrator`
	 Orchestrator Credentials
	    - :term:`ORCHESTRATOR_USER <ORCHESTRATOR_USER>`
	    - :term:`ORCHESTRATOR_PASSWORD <ORCHESTRATOR_PASSWORD>`

   ORCHESTRATOR_USER

      Pass this option, when running your :term:`PMM Server` via
      |docker| to set the orchestrator user. You only need this
      parameter (along with :term:`ORCHESTRATOR_PASSWORD
      <ORCHESTRATOR_PASSWORD>` if you have set up a custom
      |orchestrator| user.

      This option has no effect if the
      :term:`ORCHESTRATOR_ENABLED <ORCHESTRATOR_ENABLED>` option is
      set to **false**.

      .. include:: .res/code/docker.run.orchestrator-enabled.orchestrator-user.orchestrator-password.txt

   ORCHESTRATOR_PASSWORD

      Pass this option, when running your :term:`PMM Server` via |docker| to set
      the orchestrator password.

      This option has no effect if the
      :term:`ORCHESTRATOR_ENABLED <ORCHESTRATOR_ENABLED>`
      option is set to **false**.

      .. include:: .res/code/docker.run.orchestrator-enabled.orchestrator-user.orchestrator-password.txt

      .. seealso:: :term:`ORCHESTRATOR_ENABLED <ORCHESTRATOR_ENABLED>`

   SERVER_USER

      By default, the user name is ``pmm``. Use this option to use another user
      name.

      |tip.run-this.root|.

      .. include:: .res/code/docker.run.server-user.txt

   SERVER_PASSWORD

      Set the password to access the |pmm-server| web interface.

      |tip.run-this.root|.

      .. include:: .res/code/docker.run.server-password.txt
      
      By default, the user name is ``pmm``. You can change it by passing the
      :term:`SERVER_USER <SERVER_USER>` variable.

   METRICS_RESOLUTION

      This environment variable sets the minimum resolution for checking
      metrics. You should set it if the latency is higher than 1 second.

      |tip.run-this.root|.

      .. include:: .res/code/docker.run.metrics-resolution.txt

   METRICS_MEMORY

      By default, |prometheus| in |pmm-server| uses all available memory for
      storing the most recently used data chunks.  Depending on the amount of
      data coming into |prometheus|, you may require to allow less memory
      consumption if it is needed for other processes.

      .. include:: .res/contents/important.option.metrics-memory.txt

      If you are still using a version of |pmm| prior to 1.13 you might need to
      set the metrics memory by passing the |opt.metrics-memory| environment
      variable along with the |docker.run| command.

      |tip.run-this.root|. The value must be passed in kilobytes. For example,
      to set the limit to 4 GB of memory run the following command:

      .. include:: .res/code/docker.run.metrics-memory.txt

      .. seealso:: 

	 |docker| documentation: Controlling memory usage in a |docker| container
	    https://docs.docker.com/config/containers/resource_constraints/

   DISABLE_UPDATES

      To update your |pmm| from web interface you only need to click the
      |gui.update| on the home page. The |opt.disable-updates| option is useful
      if updating is not desirable. Set it to **true** when running |pmm| in
      the |docker| container.

      |tip.run-this.root|.

      .. include:: .res/code/docker.run.disable-updates.txt

      The |opt.disable-updates| option removes the |gui.update| button
      from the interface and prevents the system from being updated manually.

.. _v.percona.com: http://v.percona.com

.. include:: .res/replace.txt
