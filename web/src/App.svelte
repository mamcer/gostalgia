<script lang="ts">
  import Omnisearch from './lib/Omnisearch.svelte';
  import Sidebar from './lib/Sidebar.svelte';
  import DirectoryView from './lib/DirectoryView.svelte';
  import './app.css';

  let currentView = 'home';
  let currentId: number | null = null;

  function handleNavigate(event: CustomEvent) {
    currentView = event.detail.view;
    currentId = event.detail.id || null;
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
</script>

<div class="flex min-h-screen bg-gray-950 text-white">
  <Sidebar on:navigate={handleNavigate} />

  <main class="flex-1 ml-64 p-8">
    <div class="max-w-6xl mx-auto space-y-12">
      
      <header class="flex items-center justify-between gap-8">
        <div class="flex-1">
          <Omnisearch on:navigate={handleNavigate} />
        </div>
        <div class="flex items-center gap-4">
          <button class="p-2 text-gray-400 hover:text-white transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
        </div>
      </header>

      {#if currentView === 'home'}
        <div class="flex flex-col items-center gap-12 mt-10">
          <div class="text-center space-y-4">
            <h1 class="text-6xl font-bold tracking-tighter bg-gradient-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
              Nostalgia
            </h1>
            <p class="text-gray-400 text-lg">Your intelligent media universe explorer</p>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 w-full mt-12">
            <div class="p-6 bg-gray-900/50 border border-gray-800 rounded-2xl hover:border-blue-500/50 transition-colors cursor-pointer group">
              <h3 class="text-xl font-semibold mb-2 group-hover:text-blue-400">Latest Scans</h3>
              <p class="text-gray-500 text-sm">Quick access to your most recent indexed directories.</p>
            </div>
            <div class="p-6 bg-gray-900/50 border border-gray-800 rounded-2xl hover:border-purple-500/50 transition-colors cursor-pointer group">
              <h3 class="text-xl font-semibold mb-2 group-hover:text-purple-400">Popular Tags</h3>
              <p class="text-gray-500 text-sm">Explore your collection through frequently used labels.</p>
            </div>
            <div class="p-6 bg-gray-900/50 border border-gray-800 rounded-2xl hover:border-green-500/50 transition-colors cursor-pointer group">
              <h3 class="text-xl font-semibold mb-2 group-hover:text-green-400">Media Stats</h3>
              <p class="text-gray-500 text-sm">Overview of your storage and file distribution.</p>
            </div>
          </div>
        </div>
      {:else if currentView === 'directory' && currentId}
        <DirectoryView directoryId={currentId} on:navigate={handleNavigate} />
      {/if}

    </div>
  </main>
</div>

<style>
  :global(body) {
    background-color: #030712;
    margin: 0;
    display: block;
  }
</style>
