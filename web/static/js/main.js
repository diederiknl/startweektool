// web/static/js/main.js
document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.getElementById('searchInput');
    const searchResults = document.getElementById('searchResults');

    let debounceTimeout;

    searchInput.addEventListener('input', function () {
        clearTimeout(debounceTimeout);
        debounceTimeout = setTimeout(() => {
            const searchTerm = this.value;
            if (searchTerm.length >= 2) {
                fetch(`/api/search?q=${encodeURIComponent(searchTerm)}`)
                    .then(response => response.json())
                    .then(data => {
                        displayResults(data);
                    })
                    .catch(error => {
                        console.error('Error:', error);
                        searchResults.innerHTML = '<div class="search-result-item error">Er is een fout opgetreden bij het zoeken.</div>';
                    });
            } else {
                searchResults.innerHTML = '';
            }
        }, 300);
    });

    function displayResults(results) {
        searchResults.innerHTML = '';

        if (!Array.isArray(results)) {
            // Als we een foutmelding krijgen
            searchResults.innerHTML = '<div class="search-result-item error">Er is een fout opgetreden bij het zoeken.</div>';
            return;
        }

        if (results.length === 0) {
            searchResults.innerHTML = '<div class="search-result-item">Geen resultaten gevonden</div>';
            return;
        }

        results.forEach(student => {
            const div = document.createElement('div');
            div.className = 'search-result-item';
            div.innerHTML = `
                <div class="student-info">
                    <strong>${student.student_number}</strong> - 
                    ${student.first_name} ${student.last_name}
                </div>
                <div class="location-info">
                    Lokaal: <strong>${student.classroom}</strong> | 
                    Coach: ${student.coach}
                </div>
            `;
            searchResults.appendChild(div);
        });
    }
});