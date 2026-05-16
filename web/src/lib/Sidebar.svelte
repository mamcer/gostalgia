<script lang="ts">
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let recentScans = [];
  let popularTags = [];

  onMount(async () => {
    try {
      const [scansRes, tagsRes] = await Promise.all([
        fetch('http://localhost:8080/v1/scans/recent?limit=5'),
        fetch('http://localhost:8080/v1/tags/popular?limit=8')
      ]);
      
      if (scansRes.ok) recentScans = await scansRes.json();
      if (tagsRes.ok) popularTags = await tagsRes.json();
    } catch (e) {
      console.error('Failed to fetch sidebar data', e);
    }
  });
</script>

<aside class="w-64 h-screen bg-gray-900 border-r border-gray-800 flex flex-col fixed left-0 top-0">
  <div class="p-6 border-b border-gray-800">
    <h2 class="text-2xl font-bold text-white tracking-tighter cursor-pointer" on:click={() => dispatch('navigate', { view: 'home' })}>
      Nostalgia
    </h2>
  </div>

  <nav class="flex-1 overflow-y-auto p-4 space-y-8">
    
    <div>
      <h3 class="px-2 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">Library</h3>
      <div class="space-y-1">
        <button 
          on:click={() => dispatch('navigate', { view: 'home' })}
          class="w-full text-left px-3 py-2 rounded-lg text-gray-300 hover:bg-gray-800 hover:text-white transition-colors flex items-center gap-3"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
            <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
          </svg>
          Home
        </button>
      </div>
    </div>

    {#if recentScans.length > 0}
      <div>
        <h3 class="px-2 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">Recent Scans</h3>
        <div class="space-y-1">
          {#each recentScans as scan}
            <button 
              on:click={() => dispatch('navigate', { view: 'directory', id: scan.root_directory_id })}
              class="w-full text-left px-3 py-2 rounded-lg text-gray-400 hover:bg-gray-800 hover:text-white transition-colors text-sm truncate"
            >
              #{scan.id} • {new Date(scan.date_created).toLocaleDateString()}
            </button>
          {/each}
        </div>
      </div>
    {/if}

    {#if popularTags.length > 0}
      <div>
        <h3 class="px-2 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">Popular Tags</h3>
        <div class="flex flex-wrap gap-2 px-2">
          {#each popularTags as tag}
            <button class="px-2 py-1 bg-gray-800 hover:bg-blue-900/30 text-gray-400 hover:text-blue-400 border border-gray-700 rounded text-xs transition-colors">
              {tag}
            </button>
          {/each}
        </div>
      </div>
    {/if}

  </nav>

  <div class="p-4 border-t border-gray-800 bg-gray-900/50">
    <div class="flex items-center gap-3 px-2 py-2">
      <div class="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-500 to-purple-600"></div>
      <div class="flex flex-col">
        <span class="text-sm font-medium text-gray-200">Mario</span>
        <span class="text-xs text-gray-500">Admin</span>
      </div>
    </div>
  </div>
</aside>
