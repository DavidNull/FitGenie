import 'dart:io';
import 'dart:typed_data';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:image_picker/image_picker.dart';
import 'package:path_provider/path_provider.dart';
import 'package:provider/provider.dart';
import '../providers/app_provider.dart';
import '../models/clothing_item.dart';

class CameraScreen extends StatefulWidget {
  const CameraScreen({super.key});

  @override
  State<CameraScreen> createState() => _CameraScreenState();
}

class _CameraScreenState extends State<CameraScreen> {
  File? _selectedImage;
  String? _selectedLocalAsset;
  String? _selectedCategory;
  bool _isLoading = false;
  
  // Imágenes locales disponibles para testing
  final List<Map<String, String>> _localAssets = [
    {'path': 'assets/clothing/c1.jpg', 'name': 'Camiseta Azul'},
    {'path': 'assets/clothing/c2.jpg', 'name': 'Camisa Blanca'},
    {'path': 'assets/clothing/p1.jpg', 'name': 'Pantalón Negro'},
    {'path': 'assets/clothing/p2.jpg', 'name': 'Zapatillas'},
  ];

  final List<String> _topCategories = ['Camisetas', 'Camisas', 'Chaquetas', 'Sudaderas'];
  final List<String> _bottomCategories = ['Pantalones', 'Vaqueros', 'Shorts', 'Faldas'];
  final List<String> _shoeCategories = ['Calzado', 'Zapatillas', 'Zapatos'];

  Future<void> _pickImage(ImageSource source) async {
    final picker = ImagePicker();
    final pickedFile = await picker.pickImage(source: source);
    
    if (pickedFile != null) {
      setState(() {
        _selectedImage = File(pickedFile.path);
        _selectedLocalAsset = null; // Limpiar asset local
      });
    }
  }
  
  void _selectLocalAsset(String assetPath) {
    setState(() {
      _selectedLocalAsset = assetPath;
      _selectedImage = null; // Limpiar imagen de archivo
    });
  }

  Future<void> _saveClothingItem() async {
    if ((_selectedImage == null && _selectedLocalAsset == null) || _selectedCategory == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Selecciona una imagen y categoría'),
          backgroundColor: Colors.orange,
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    try {
      final provider = Provider.of<AppProvider>(context, listen: false);
      
      // Subir imagen al backend (ya sea archivo o asset local)
      String imageUrl;
      if (_selectedLocalAsset != null) {
        // Convertir asset a archivo temporal y subir
        final tempFile = await _assetToFile(_selectedLocalAsset!);
        final uploadedUrl = await provider.uploadImage(tempFile);
        if (uploadedUrl == null || uploadedUrl.isEmpty) {
          throw Exception('Failed to upload local asset');
        }
        imageUrl = uploadedUrl;
        // Limpiar archivo temporal
        await tempFile.delete();
      } else {
        final uploadedUrl = await provider.uploadImage(_selectedImage!);
        if (uploadedUrl == null || uploadedUrl.isEmpty) {
          throw Exception('Failed to upload image');
        }
        imageUrl = uploadedUrl;
      }

      final newItem = ClothingItem(
        id: '',
        userId: provider.userId ?? '',
        name: '$_selectedCategory ${DateTime.now().millisecondsSinceEpoch}',
        category: _selectedCategory!,
        imageUrl: imageUrl,
      );

      await provider.addClothingItem(newItem);
      
      // Recargar lista de prendas para que aparezca en galería
      await provider.loadClothingItems();

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Prenda guardada exitosamente'),
            backgroundColor: Color(0xFF1DA9B6),
          ),
        );
        setState(() {
          _selectedImage = null;
          _selectedLocalAsset = null;
          _selectedCategory = null;
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F8FA),
      body: SafeArea(
        child: Column(
          children: [
            Container(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  IconButton(
                    onPressed: () => Navigator.pop(context),
                    icon: const Icon(Icons.arrow_back, color: Color(0xFF0E4A88), size: 28),
                  ),
                  const Expanded(
                    child: Center(
                      child: Text(
                        'Añadir prenda',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Color(0xFF0E4A88),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 48), // Balance for back button
                ],
              ),
            ),
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                child: Column(
                  children: [
                    const Text(
                      'Añadir prenda',
                      style: TextStyle(
                        fontSize: 28,
                        fontWeight: FontWeight.bold,
                        color: Color(0xFF0E4A88),
                      ),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      'Selecciona el tipo de prenda y sube una foto',
                      style: TextStyle(
                        fontSize: 16,
                        color: Color(0xFF1DA9B6),
                      ),
                    ),
                    const SizedBox(height: 20),
                    _buildCategorySelector(),
                    const SizedBox(height: 30),
                    _buildImagePreview(),
                    const SizedBox(height: 30),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        GestureDetector(
                          onTap: () => _pickImage(ImageSource.camera),
                          child: Container(
                            width: 70,
                            height: 70,
                            decoration: BoxDecoration(
                              color: const Color(0xFFF78400),
                              shape: BoxShape.circle,
                              border: Border.all(
                                color: const Color(0xFFE9ECF1),
                                width: 4,
                              ),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.black.withOpacity(0.3),
                                  blurRadius: 10,
                                  offset: const Offset(0, 5),
                                ),
                              ],
                            ),
                            child: const Icon(
                              Icons.camera_alt,
                              size: 32,
                              color: Colors.white,
                            ),
                          ),
                        ),
                        const SizedBox(width: 20),
                        GestureDetector(
                          onTap: () => _pickImage(ImageSource.gallery),
                          child: Container(
                            width: 70,
                            height: 70,
                            decoration: BoxDecoration(
                              color: const Color(0xFF1DA9B6),
                              shape: BoxShape.circle,
                              border: Border.all(
                                color: const Color(0xFFE9ECF1),
                                width: 4,
                              ),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.black.withOpacity(0.3),
                                  blurRadius: 10,
                                  offset: const Offset(0, 5),
                                ),
                              ],
                            ),
                            child: const Icon(
                              Icons.photo_library,
                              size: 32,
                              color: Colors.white,
                            ),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 30),
                    if ((_selectedImage != null || _selectedLocalAsset != null) && _selectedCategory != null)
                      SizedBox(
                        width: double.infinity,
                        child: ElevatedButton.icon(
                          onPressed: _isLoading ? null : _saveClothingItem,
                          icon: _isLoading
                              ? const SizedBox(
                                  width: 20,
                                  height: 20,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 2,
                                    color: Colors.white,
                                  ),
                                )
                              : const Icon(Icons.save),
                          label: Text(_isLoading ? 'Guardando...' : 'Guardar prenda'),
                          style: ElevatedButton.styleFrom(
                            backgroundColor: const Color(0xFF0E4A88),
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(vertical: 16),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                        ),
                      ),
                    const SizedBox(height: 30),
                    const Divider(),
                    const SizedBox(height: 16),
                    const Text(
                      'O usa una imagen local (modo desarrollo):',
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.grey,
                      ),
                    ),
                    const SizedBox(height: 12),
                    _buildLocalAssetsSelector(),
                    const SizedBox(height: 20),
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: const Color(0xFF1DA9B6).withOpacity(0.2),
                        borderRadius: BorderRadius.circular(12),
                      ),
                      child: const Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            Icons.lightbulb,
                            color: Color(0xFFF78400),
                            size: 20,
                          ),
                          SizedBox(width: 8),
                          Text(
                            'Consejo: Buena luz = mejor análisis',
                            style: TextStyle(
                              color: Color(0xFF0E4A88),
                              fontSize: 14,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 20),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCategorySelector() {
    return Column(
      children: [
        _buildCategorySection('Parte de arriba', _topCategories, Icons.accessibility_new),
        const SizedBox(height: 12),
        _buildCategorySection('Parte de abajo', _bottomCategories, Icons.accessibility),
        const SizedBox(height: 12),
        _buildCategorySection('Calzado', _shoeCategories, Icons.directions_walk),
      ],
    );
  }

  Widget _buildCategorySection(String title, List<String> categories, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(icon, color: const Color(0xFF0E4A88), size: 20),
              const SizedBox(width: 8),
              Text(
                title,
                style: const TextStyle(
                  fontWeight: FontWeight.bold,
                  color: Color(0xFF0E4A88),
                  fontSize: 14,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: categories.map((category) {
              final isSelected = _selectedCategory == category;
              return GestureDetector(
                onTap: () {
                  setState(() {
                    _selectedCategory = category;
                  });
                },
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: isSelected ? const Color(0xFF0E4A88) : const Color(0xFFE9ECF1),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Text(
                    category,
                    style: TextStyle(
                      color: isSelected ? Colors.white : const Color(0xFF0E4A88),
                      fontSize: 12,
                      fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                    ),
                  ),
                ),
              );
            }).toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildImagePreview() {
    // Preview de asset local
    if (_selectedLocalAsset != null) {
      return Container(
        width: 280,
        height: 280,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(24),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.2),
              blurRadius: 20,
              offset: const Offset(0, 10),
            ),
          ],
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(24),
          child: Image.asset(
            _selectedLocalAsset!,
            fit: BoxFit.cover,
          ),
        ),
      );
    }
    
    // Preview de archivo seleccionado
    if (_selectedImage != null) {
      return Container(
        width: 280,
        height: 280,
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(24),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.2),
              blurRadius: 20,
              offset: const Offset(0, 10),
            ),
          ],
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(24),
          child: Image.file(
            _selectedImage!,
            fit: BoxFit.cover,
          ),
        ),
      );
    }

    // Placeholder vacío
    return Container(
      width: 280,
      height: 280,
      decoration: BoxDecoration(
        color: const Color(0xFFE9ECF1),
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          color: const Color(0xFF1DA9B6),
          width: 4,
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.2),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: const Color(0xFFF78400).withOpacity(0.1),
              shape: BoxShape.circle,
            ),
            child: const Icon(
              Icons.camera_alt,
              size: 64,
              color: Color(0xFFF78400),
            ),
          ),
          const SizedBox(height: 20),
          const Text(
            'Toca la cámara o galería',
            style: TextStyle(
              color: Color(0xFF0E4A88),
              fontSize: 16,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLocalAssetsSelector() {
    return SizedBox(
      height: 100,
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: _localAssets.length,
        itemBuilder: (context, index) {
          final asset = _localAssets[index];
          final isSelected = _selectedLocalAsset == asset['path'];
          
          return GestureDetector(
            onTap: () => _selectLocalAsset(asset['path']!),
            child: Container(
              width: 80,
              margin: const EdgeInsets.only(right: 12),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(12),
                border: isSelected
                    ? Border.all(color: const Color(0xFF0E4A88), width: 3)
                    : null,
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.1),
                    blurRadius: 5,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: ClipRRect(
                borderRadius: BorderRadius.circular(12),
                child: Stack(
                  fit: StackFit.expand,
                  children: [
                    Image.asset(
                      asset['path']!,
                      fit: BoxFit.cover,
                    ),
                    if (isSelected)
                      Container(
                        color: const Color(0xFF0E4A88).withOpacity(0.5),
                        child: const Icon(
                          Icons.check,
                          color: Colors.white,
                          size: 32,
                        ),
                      ),
                  ],
                ),
              ),
            ),
          );
        },
      ),
    );
  }

  // Helper: Convert asset to temporary file for upload
  Future<File> _assetToFile(String assetPath) async {
    final byteData = await rootBundle.load(assetPath);
    final bytes = byteData.buffer.asUint8List();
    
    final tempDir = await getTemporaryDirectory();
    final fileName = assetPath.split('/').last;
    final tempFile = File('${tempDir.path}/$fileName');
    
    await tempFile.writeAsBytes(bytes);
    return tempFile;
  }
}
