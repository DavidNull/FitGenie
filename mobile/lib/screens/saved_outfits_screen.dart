import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/app_provider.dart';
import '../models/outfit.dart';
import '../models/clothing_item.dart';

class SavedOutfitsScreen extends StatefulWidget {
  const SavedOutfitsScreen({super.key});

  @override
  State<SavedOutfitsScreen> createState() => _SavedOutfitsScreenState();
}

class _SavedOutfitsScreenState extends State<SavedOutfitsScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final provider = context.read<AppProvider>();
      if (provider.userId != null) {
        provider.loadOutfits(provider.userId!);
      }
    });
  }

  List<Outfit> _getFavoriteOutfits(List<Outfit> outfits) {
    return outfits.where((o) => o.favorite).toList();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFE9ECF1),
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              Color(0xFFF0F7F8),
              Color(0xFFF0F7F8),
              Color(0xFF0E4A88),
            ],
            stops: [0.0, 0.4, 1.0],
          ),
        ),
        child: SafeArea(
          child: Consumer<AppProvider>(
            builder: (context, provider, child) {
              final favoriteOutfits = _getFavoriteOutfits(provider.outfits);
              
              return SingleChildScrollView(
                child: Padding(
                  padding: const EdgeInsets.all(20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          IconButton(
                            onPressed: () => Navigator.pop(context),
                            icon: const Icon(Icons.arrow_back, color: Color(0xFF0E4A88)),
                          ),
                          const Expanded(
                            child: Text(
                              'Outfits Guardados',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 28,
                                fontWeight: FontWeight.bold,
                                color: Color(0xFF0E4A88),
                              ),
                            ),
                          ),
                          const SizedBox(width: 48),
                        ],
                      ),
                      const SizedBox(height: 8),
                      Text(
                        '${favoriteOutfits.length} outfits favoritos',
                        style: const TextStyle(
                          fontSize: 16,
                          color: Color(0xFF1DA9B6),
                        ),
                      ),
                      const SizedBox(height: 20),
                      if (provider.isLoading)
                        const Center(
                          child: Padding(
                            padding: EdgeInsets.all(40),
                            child: CircularProgressIndicator(),
                          ),
                        )
                      else if (favoriteOutfits.isEmpty)
                        Center(
                          child: Padding(
                            padding: const EdgeInsets.all(40),
                            child: Column(
                              children: [
                                Icon(
                                  Icons.favorite_border,
                                  size: 64,
                                  color: Colors.grey[400],
                                ),
                                const SizedBox(height: 16),
                                Text(
                                  'No tienes outfits favoritos.\nGuarda outfits desde las recomendaciones.',
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontSize: 16,
                                    color: Colors.grey[600],
                                  ),
                                ),
                              ],
                            ),
                          ),
                        )
                      else
                        ListView.builder(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          itemCount: favoriteOutfits.length,
                          itemBuilder: (context, index) {
                            return _buildOutfitCard(favoriteOutfits[index]);
                          },
                        ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }

  Widget _buildOutfitCard(Outfit outfit) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    outfit.name,
                    style: const TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: Color(0xFF0E4A88),
                    ),
                  ),
                ),
                IconButton(
                  onPressed: () => _removeFromFavorites(outfit),
                  icon: const Icon(Icons.favorite, color: Colors.red),
                ),
              ],
            ),
          ),
          if (outfit.clothingItems.isNotEmpty)
            SizedBox(
              height: 120,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                padding: const EdgeInsets.symmetric(horizontal: 16),
                itemCount: outfit.clothingItems.length,
                itemBuilder: (context, index) {
                  return _buildItemPreview(outfit.clothingItems[index]);
                },
              ),
            )
          else
            Padding(
              padding: const EdgeInsets.all(16),
              child: Text(
                'Sin prendas asociadas',
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.grey[600],
                ),
              ),
            ),
          if (outfit.notes != null && outfit.notes!.isNotEmpty)
            Padding(
              padding: const EdgeInsets.all(16),
              child: Text(
                outfit.notes!,
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.grey[600],
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildItemPreview(ClothingItem item) {
    return Container(
      width: 100,
      margin: const EdgeInsets.only(right: 12),
      decoration: BoxDecoration(
        color: const Color(0xFFE9ECF1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(12),
        child: item.imageUrl != null && item.imageUrl!.isNotEmpty
            ? Image.network(
                item.imageUrl!,
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) => _buildItemPlaceholder(item),
              )
            : _buildItemPlaceholder(item),
      ),
    );
  }

  Widget _buildItemPlaceholder(ClothingItem item) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            _getCategoryIcon(item.category),
            size: 32,
            color: const Color(0xFF0E4A88).withOpacity(0.5),
          ),
          const SizedBox(height: 4),
          Text(
            item.category,
            style: TextStyle(
              fontSize: 10,
              color: Colors.grey[600],
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }

  Future<void> _removeFromFavorites(Outfit outfit) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Eliminar de favoritos'),
        content: Text('¿Quitar "${outfit.name}" de tus favoritos?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancelar'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Eliminar'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      final provider = Provider.of<AppProvider>(context, listen: false);
      // Create updated outfit with favorite = false
      final updatedOutfit = Outfit(
        id: outfit.id,
        userId: outfit.userId,
        name: outfit.name,
        clothingItems: outfit.clothingItems,
        occasion: outfit.occasion,
        season: outfit.season,
        style: outfit.style,
        overallScore: outfit.overallScore,
        rating: outfit.rating,
        worn: outfit.worn,
        favorite: false,
        notes: outfit.notes,
        createdAt: outfit.createdAt,
        updatedAt: DateTime.now(),
      );
      
      await provider.updateOutfit(updatedOutfit);
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Eliminado de favoritos'),
            backgroundColor: Color(0xFF1DA9B6),
          ),
        );
      }
    }
  }

  IconData _getCategoryIcon(String category) {
    final lower = category.toLowerCase();
    if (lower.contains('top') || 
        lower.contains('shirt') || 
        lower.contains('camiseta') ||
        lower.contains('parte de arriba')) {
      return Icons.checkroom;
    }
    if (lower.contains('bottom') || 
        lower.contains('pant') || 
        lower.contains('parte de abajo')) {
      return Icons.accessibility;
    }
    if (lower.contains('shoe') || 
        lower.contains('calzado')) {
      return Icons.directions_walk;
    }
    return Icons.checkroom;
  }
}
