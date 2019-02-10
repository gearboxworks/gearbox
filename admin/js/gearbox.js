gearbox.loadProjects();
function getGearboxProjects() {
    if (!gearbox.projects) {
        window.setTimeout(getGearboxProjects,50);
        return;
    }
    var app = new Vue({
        el: '#app',
        data: {
            projects: gearbox.projects,
        },
    });
}
getGearboxProjects();