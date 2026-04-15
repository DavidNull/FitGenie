import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/app_provider.dart';
import '../models/outfit_recommendation.dart';

class RecommendationsScreen extends StatefulWidget {
  const RecommendationsScreen({super.key});

  @override
  State<RecommendationsScreen> createState() => _RecommendationsScreenState();
}

class _RecommendationsScreenState extends State<RecommendationsScreen> {
  String? _selectedOccasion;
  String? _selectedSeason;

  final List<String> _occasions = [
    'Casual',
    'Formal',
    'Trabajo',
    'Fiesta',
    'Deporte',
  ];

  final List<String> _seasons = [
    'Primavera',
    'Verano',
    'Otoño',
    'Invierno',
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _loadRecommendations();
    });
  }

  Future<void> _loadRecommendations() async {
    final provider = Provider.of<AppProvider>(context, listen: false);
    await provider.getRecommendations(
      occasion: _selectedOccasion,
      season: _selectedSeason,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: Color(0xFF0E4A88)),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: const Text(
          'Recomendaciones IA',
          style: TextStyle(
            color: Color(0xFF0E4A88),
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
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
              return Column(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          '${provider.recommendations.length} outfits generados',
                          style: const TextStyle(
                            fontSize: 16,
                            color: Color(0xFF1DA9B6),
                          ),
                        ),
                        const SizedBox(height: 16),
                        Row(
                          children: [
                            Expanded(
                              child: _buildDropdown(
                                'Ocasión',
                                _selectedOccasion,
                                _occasions,
                                (value) {
                                  setState(() => _selectedOccasion = value);
                                  _loadRecommendations();
                                },
                              ),
                            ),
                            const SizedBox(width: 12),
                            Expanded(
                              child: _buildDropdown(
                                'Temporada',
                                _selectedSeason,
                                _seasons,
                                (value) {
                                  setState(() => _selectedSeason = value);
                                  _loadRecommendations();
                                },
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  Expanded(
                    child: provider.recommendations.isEmpty
                        ? _buildEmptyState()
                        : ListView.builder(
                            padding: const EdgeInsets.symmetric(horizontal: 20),
                            itemCount: provider.recommendations.length,
                            itemBuilder: (context, index) {
                              return _buildRecommendationCard(
                                provider.recommendations[index],
                              );
                            },
                          ),
                  ),
                ],
              );
            },
          ),
        ),
      ),
    );
  }

  Widget _buildDropdown(
    String label,
    String? value,
    List<String> items,
    Function(String?) onChanged,
  ) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12),
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
      child: DropdownButtonHideUnderline(
        child: DropdownButton<String>(
          value: value,
          hint: Text(label),
          isExpanded: true,
          items: [
            const DropdownMenuItem(
              value: null,
              child: Text('Todas'),
            ),
            ...items.map((item) => DropdownMenuItem(
              value: item,
              child: Text(item),
            )),
          ],
          onChanged: onChanged,
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.style_outlined,
            size: 64,
            color: Color(0xFF1DA9B6),
          ),
          SizedBox(height: 16),
          Text(
            'No hay recomendaciones',
            style: TextStyle(
              fontSize: 18,
              color: Color(0xFF0E4A88),
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: 8),
          Text(
            'Agrega más prendas para generar outfits',
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRecommendationCard(OutfitRecommendation rec) {
    final outfit = rec.outfit;
    if (outfit == null) return const SizedBox.shrink();

    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.08),
            blurRadius: 12,
            offset: const Offset(0, 4),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            decoration: const BoxDecoration(
              color: Color(0xFF0E4A88),
              borderRadius: BorderRadius.vertical(
                top: Radius.circular(20),
              ),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    outfit.name,
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 6,
                  ),
                  decoration: BoxDecoration(
                    color: Colors.white.withOpacity(0.2),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Text(
                    '${(rec.confidence * 100).toInt()}%',
                    style: const TextStyle(
                      color: Colors.white,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  rec.reasoning,
                  style: const TextStyle(
                    fontSize: 14,
                    color: Colors.grey,
                    fontStyle: FontStyle.italic,
                  ),
                ),
                const SizedBox(height: 12),
                const Text(
                  'Prendas:',
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF0E4A88),
                  ),
                ),
                const SizedBox(height: 8),
                ...outfit.clothingItems.map((item) => Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Row(
                    children: [
                      const Icon(
                        Icons.checkroom,
                        size: 16,
                        color: Color(0xFF1DA9B6),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        '${item.name} (${item.category})',
                        style: const TextStyle(fontSize: 14),
                      ),
                    ],
                  ),
                )),
                const SizedBox(height: 12),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton.icon(
                    onPressed: () => _saveOutfit(rec),
                    icon: const Icon(Icons.save),
                    label: const Text('Guardar Outfit'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF0E4A88),
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 12),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _saveOutfit(OutfitRecommendation rec) async {
    final provider = Provider.of<AppProvider>(context, listen: false);
    await provider.createOutfitFromRecommendation(rec);
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Outfit guardado exitosamente'),
          backgroundColor: Color(0xFF1DA9B6),
        ),
      );
    }
  }
}
