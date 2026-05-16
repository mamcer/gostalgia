<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import PreviewModal from './PreviewModal.svelte';

  const dispatch = createEventDispatcher();

  export let directoryId: number;

  let directory = null;
  let loading = true;
  let viewMode: 'grid' | 'list' = 'grid';
  
  // Filters
  let filterType = 'all';
  let filterSize = 0; // min size in KB
  let searchQuery = '';

  // Preview
  let selectedFile = null;
  let showPreview = false;

  const API_URL = 'http://localhost:8080';

  async function fetchDirectory(id: number) {
    loading = true;
    try {
      const response = await fetch(`${API_URL}/v1/directories/${id}`);
      if (response.ok) {
        directory = await response.json();
      }
    } catch (e) {
      console.error('Failed to fetch directory', e);
    } finally {
      // Small delay for skeleton visibility
      setTimeout(() => { loading = false; }, 300);
    }
  }

  function isImage(extension: string) {
    const images = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp'];
    return images.includes(extension.toLowerCase());
  }

  function openPreview(file: any) {
    selectedFile = file;
    showPreview = true;
  }

  $: filteredFiles = directory?.files?.filter(f => {
    const matchesType = filterType === 'all' || isOfType(f.extension, filterType);
    const matchesSize = f.size_raw ? (f.size_raw / 1024) >= filterSize : true;
    const matchesSearch = f.name.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesType && matchesSize && matchesSearch;
  }) || [];

  function isOfType(ext: string, type: string) {
    const types = {
      image: ['.jpg', '.jpeg', '.png', '.gif', '.webp'],
      video: ['.mp4', '.mkv', '.avi', '.mov'],
      doc: ['.pdf', '.doc', '.docx', '.txt', '.md'],
      zip: ['.zip', '.rar', '.7z', '.tar', '.gz']
    };
    return types[type]?.includes(ext.toLowerCase());
  }

  $: if (directoryId) fetchDirectory(directoryId);

  onMount(() => fetchDirectory(directoryId));
</script>

<div class="space-y-6">
  
  {#if directory || loading}
    <!-- Breadcrumbs -->
    <div class="flex items-center gap-2 text-sm text-gray-500 overflow-x-auto whitespace-nowrap pb-2">
      <button on:click={() => dispatch('navigate', { view: 'home' })} class="hover:text-white transition-colors">Root</button>
      <span>/</span>
      <span class="text-gray-300 font-medium">{directory?.name || '...'}</span>
    </div>

    <!-- Header & Toggle -->
    <div class="flex items-end justify-between border-b border-gray-800 pb-6">
      <div>
        <h1 class="text-3xl font-bold text-white">{directory?.name || 'Loading...'}</h1>
        <p class="text-gray-500 text-sm mt-1">{directory?.full_path || 'Please wait'}</p>
      </div>
      
      <div class="flex items-center gap-2 bg-gray-900 p-1 rounded-lg border border-gray-800">
        <button 
          on:click={() => viewMode = 'grid'}
          class="p-2 rounded-md transition-all {viewMode === 'grid' ? 'bg-gray-800 text-blue-400' : 'text-gray-500 hover:text-gray-300'}"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v6a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM11 13a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </button>
        <button 
          on:click={() => viewMode = 'list'}
          class="p-2 rounded-md transition-all {viewMode === 'list' ? 'bg-gray-800 text-blue-400' : 'text-gray-500 hover:text-gray-300'}"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h14a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h14a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h14a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h14a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="flex flex-wrap items-center gap-4 bg-gray-900/50 p-4 rounded-xl border border-gray-800">
      <div class="flex-1 min-w-[200px] relative">
        <input 
          type="text" 
          bind:value={searchQuery}
          placeholder="Filter current view..." 
          class="w-full bg-gray-800 border-gray-700 text-sm rounded-lg pl-10 focus:ring-blue-500"
        />
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 absolute left-3 top-3 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
        </svg>
      </div>

      <select bind:value={filterType} class="bg-gray-800 border-gray-700 text-sm rounded-lg focus:ring-blue-500">
        <option value="all">All Types</option>
        <option value="image">Images</option>
        <option value="video">Videos</option>
        <option value="doc">Documents</option>
        <option value="zip">Archives</option>
      </select>

      <div class="flex items-center gap-2">
        <span class="text-xs text-gray-500 font-medium whitespace-nowrap">Min Size:</span>
        <select bind:value={filterSize} class="bg-gray-800 border-gray-700 text-sm rounded-lg focus:ring-blue-500">
          <option value={0}>Any</option>
          <option value={1024}>1 MB</option>
          <option value={10240}>10 MB</option>
          <option value={102400}>100 MB</option>
          <option value={1048576}>1 GB</option>
        </select>
      </div>

      <div class="ml-auto text-xs text-gray-500">
        Showing <b>{filteredFiles.length}</b> of {directory?.files?.length || 0} files
      </div>
    </div>

    {#if loading}
      <!-- SKELETON GRID -->
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
        {#each Array(16) as _}
          <div class="animate-pulse bg-gray-900/50 border border-gray-800 rounded-xl p-4 space-y-3">
            <div class="aspect-square bg-gray-800 rounded-lg"></div>
            <div class="h-2 bg-gray-800 rounded w-3/4 mx-auto"></div>
            <div class="h-2 bg-gray-800 rounded w-1/2 mx-auto"></div>
          </div>
        {/each}
      </div>
    {:else}
      
      {#if viewMode === 'grid'}
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
          {#if directory.directories}
            {#each directory.directories as sub}
              <button 
                on:click={() => dispatch('navigate', { view: 'directory', id: sub.id })}
                class="flex flex-col items-center p-4 rounded-xl bg-gray-900/30 border border-gray-800 hover:border-yellow-500/50 hover:bg-gray-800/50 transition-all group"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-yellow-500/80 mb-2 group-hover:scale-110 transition-transform" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                </svg>
                <span class="text-xs text-gray-300 font-medium text-center line-clamp-2">{sub.name}</span>
              </button>
            {/each}
          {/if}

          {#if filteredFiles}
            {#each filteredFiles as file}
              <div 
                on:click={() => openPreview(file)}
                class="flex flex-col items-center p-2 rounded-xl bg-gray-900/20 border border-gray-800/50 hover:border-blue-500/30 transition-all group cursor-pointer"
              >
                <div class="relative w-full aspect-square mb-2 bg-gray-800 rounded-lg flex items-center justify-center overflow-hidden">
                  {#if isImage(file.extension)}
                    <img 
                      src={`${API_URL}/thumbs/${file.path}`} 
                      alt={file.name}
                      loading="lazy"
                      class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                      on:error={(e) => e.target.src = '/icons.svg'} 
                    />
                  {:else}
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-600" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd" />
                    </svg>
                  {/if}
                </div>
                <span class="text-[10px] text-gray-400 font-medium text-center line-clamp-1 w-full px-1">{file.name}</span>
                <span class="text-[9px] text-gray-600 mt-0.5">{file.size}</span>
              </div>
            {/each}
          {/if}
        </div>
      {:else}
        <!-- LIST VIEW -->
        <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden">
          <table class="w-full text-left text-sm">
            <thead class="bg-gray-900/50 text-gray-500 font-medium border-b border-gray-800">
              <tr>
                <th class="px-4 py-3">Name</th>
                <th class="px-4 py-3">Date</th>
                <th class="px-4 py-3">Size</th>
                <th class="px-4 py-3">Type</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-800/50">
              {#if directory.directories}
                {#each directory.directories as sub}
                  <tr 
                    on:click={() => dispatch('navigate', { view: 'directory', id: sub.id })}
                    class="hover:bg-gray-800/30 cursor-pointer transition-colors group"
                  >
                    <td class="px-4 py-3 flex items-center gap-3">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-500/70" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                      </svg>
                      <span class="text-gray-200 group-hover:text-white font-medium">{sub.name}</span>
                    </td>
                    <td class="px-4 py-3 text-gray-500">{sub.date_modified}</td>
                    <td class="px-4 py-3 text-gray-500">{sub.size}</td>
                    <td class="px-4 py-3 text-gray-600">Folder</td>
                  </tr>
                {/each}
              {/if}
              {#if filteredFiles}
                {#each filteredFiles as file}
                  <tr 
                    on:click={() => openPreview(file)}
                    class="hover:bg-gray-800/20 transition-colors group cursor-pointer"
                  >
                    <td class="px-4 py-3 flex items-center gap-3">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd" />
                      </svg>
                      <span class="text-gray-300 group-hover:text-gray-100">{file.name}</span>
                    </td>
                    <td class="px-4 py-3 text-gray-500">{file.date_modified}</td>
                    <td class="px-4 py-3 text-gray-500">{file.size}</td>
                    <td class="px-4 py-3 text-gray-600 uppercase text-[10px]">{file.extension.replace('.','')}</td>
                  </tr>
                {/each}
              {/if}
            </tbody>
          </table>
        </div>
      {/if}

    {/if}
  {/if}
</div>

<PreviewModal 
  file={selectedFile} 
  isOpen={showPreview} 
  on:close={() => showPreview = false} 
/>
