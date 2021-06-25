
import * as THREE from 'https://cdn.skypack.dev/three@0.129.0';

const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
const renderer = new THREE.WebGLRenderer();
const scene = new THREE.Scene();

renderer.setSize(window.innerWidth, window.innerHeight);
document.body.appendChild(renderer.domElement);

const size = 1000;
const divisions = 500;

const gridHelper = new THREE.GridHelper(size, divisions, 0x7a04eb, 0x7a04eb);

scene.add(gridHelper);

camera.position.y = 5;
const animate = function () {
    camera.position.z -= 0.05;
    requestAnimationFrame( animate );
    renderer.render( scene, camera );
};

animate();

window.addEventListener('resize', function(event) {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
})