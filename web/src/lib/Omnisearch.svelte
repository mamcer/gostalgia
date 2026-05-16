<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let query = '';
  let results = { files: [], directories: [], tags: [], total: 0 };
  let loading = false;
  let showResults = false;
  let timer: ReturnType<typeof setTimeout>;

  async function search() {
    if (query.length < 2) {
      results = { files: [], directories: [], tags: [], total: 0 };
      showResults = false;
      return;
    }

    loading = true;
    try {
      const response = await fetch(`http://localhost:8080/v1/search?q=${encodeURIComponent(query)}`);
      if (response.ok) {
        results = await response.json();
        showResults = true;
      }
    } catch (error) {
      console.error('Search error:', error);
    } finally {
      loading = false;
    }
  }

  function handleInput() {
    clearTimeout(timer);
    timer = setTimeout(search, 300);
  }

  function closeSearch() {
    setTimeout(() => { showResults = false; }, 200);
  }

  function selectDirectory(id: number) {
    dispatch('navigate', { view: 'directory', id });
    showResults = false;
    query = '';
  }
</script>

<div class="relative w-full max-w-2xl mx-auto">
  <div class="relative">
    <input
      type="text"
      bind:value={query}
      on:input={handleInput}
      on:blur={closeSearch}
      on:focus={() => query.length >= 2 && (showResults = true)}
      placeholder="Search anything (files, folders, tags)..."
      class="w-full px-4 py-3 pl-12 text-white bg-gray-800 border border-gray-700 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all shadow-lg"
    />
    <div class="absolute left-4 top-3.5 text-gray-400">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>
    {#if loading}
      <div class="absolute right-4 top-3.5">
        <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-500"></div>
      </div>
    {/if}
  </div>

  {#if showResults && (results.tags.length > 0 || results.directories.length > 0 || results.files.length > 0)}
    <div class="absolute z-50 w-full mt-2 bg-gray-900 border border-gray-700 rounded-xl shadow-2xl overflow-hidden max-h-[80vh] overflow-y-auto">
      
      {#if results.tags.length > 0}
        <div class="p-2">
          <h3 class="px-3 py-1 text-xs font-semibold text-gray-500 uppercase tracking-wider">Tags</h3>
          {#each results.tags as tag}
            <button class="w-full text-left px-3 py-2 rounded-lg hover:bg-gray-800 transition-colors flex items-center gap-2">
              <span class="text-blue-400">#</span>
              <span class="text-gray-200">{tag}</span>
            </button>
          {/each}
        </div>
      {/if}

      {#if results.directories.length > 0}
        <div class="p-2 border-t border-gray-800">
          <h3 class="px-3 py-1 text-xs font-semibold text-gray-500 uppercase tracking-wider">Directories</h3>
          {#each results.directories as dir}
            <button 
              on:click={() => selectDirectory(dir.id)}
              class="w-full text-left px-3 py-2 rounded-lg hover:bg-gray-800 transition-colors flex items-center gap-2 group"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-500" viewBox="0 0 20 20" fill="currentColor">
                <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
              </svg>
              <div class="flex flex-col">
                <span class="text-gray-200 font-medium">{dir.name}</span>
                <span class="text-xs text-gray-500 truncate">{dir.full_path}</span>
              </div>
            </button>
          {/each}
        </div>
      {/if}

      {#if results.files.length > 0}
        <div class="p-2 border-t border-gray-800">
          <h3 class="px-3 py-1 text-xs font-semibold text-gray-500 uppercase tracking-wider">Files</h3>
          {#each results.files as file}
            <button class="w-full text-left px-3 py-2 rounded-lg hover:bg-gray-800 transition-colors flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd" />
              </svg>
              <div class="flex flex-col">
                <span class="text-gray-200 font-medium">{file.name}</span>
                <span class="text-xs text-gray-500">{file.size} • {file.date_modified}</span>
              </div>
            </button>
          {/each}
        </div>
      {/if}

      <div class="p-3 bg-gray-800/50 border-t border-gray-800 text-center">
        <span class="text-xs text-gray-400">Found {results.total} results</span>
      </div>
    </div>
  {:else if showResults && query.length >= 2}
    <div class="absolute z-50 w-full mt-2 p-8 bg-gray-900 border border-gray-700 rounded-xl shadow-2xl text-center">
      <p class="text-gray-400 text-sm">No results found for "{query}"</p>
    </div>
  {/if}
</div>
